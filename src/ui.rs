use chrono::prelude::*;
use crossterm::{
    event::{self, Event as CEvent, KeyCode, KeyEvent},
    terminal::{disable_raw_mode, enable_raw_mode},
};
use serde::{Deserialize, Serialize};
use std::{fs, sync::mpsc::Receiver, fmt::Display};
use std::io;
use std::sync::mpsc;
use std::thread;
use std::time::{Duration, Instant};
use tui::{
    backend::{CrosstermBackend, Backend},
    layout::{Alignment, Constraint, Direction, Layout, Rect},
    style::{Color, Modifier, Style},
    text::{Span, Spans},
    widgets::{
        Block, BorderType, Borders, Cell, List, ListItem, ListState, Paragraph, Row, Table, Tabs,
    },
    Terminal, Frame, terminal,
};

enum Event<T> {
    Input(T),
    Tick,
}

enum HandleEventMsg {
    Continue,
    Exit
}

#[derive(Clone)]
enum ViewState {
    Main,
    Workflows
}

impl ViewState {
    fn next(&self) -> ViewState {
        match self {
            ViewState::Main => ViewState::Workflows,
            ViewState::Workflows => ViewState::Main
        }
    }
}


impl Display for ViewState {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            ViewState::Main => write!(f, "Main"),
            ViewState::Workflows => write!(f, "Workflows"),
        }
    }
}
struct App<'a> {
    view_state: ViewState,
    terminal: &'a mut Terminal<CrosstermBackend<io::Stdout>>,
}

impl App<'_> {
    fn new(terminal: &mut Terminal<CrosstermBackend<io::Stdout>>) -> App {
        App {
            view_state: ViewState::Main,
            terminal: terminal
        }
    }

    fn handle_event<'a>(self: &mut Self, rx: &'a Receiver<Event<KeyEvent>>) -> Result<HandleEventMsg, std::io::Error> {
        match rx.recv().expect("should be an unwrappable event") {
            Event::Input(event) => match event.code {
                KeyCode::Char('q') => {
                    disable_raw_mode()?;
                    self.terminal.show_cursor().expect("should show cursor");
                    Ok(HandleEventMsg::Exit)
                },
                KeyCode::Tab => {
                    self.view_state = self.view_state.next();
                    Ok(HandleEventMsg::Continue)
                },
                _ => {Ok(HandleEventMsg::Continue)}
            },
            Event::Tick => {Ok(HandleEventMsg::Continue)}
        }
    }


    fn draw(self: &mut Self, input: &str) {
        

        let view_state = self.view_state.clone();

        self.terminal.draw(move |rect| {
            let size = rect.size();
            let chunks = Layout::default()
                .direction(Direction::Vertical)
                // .margin(2)
                .constraints(
                    [
                        Constraint::Min(2),
                        Constraint::Length(3),
                    ]
                    .as_ref(),
                )
                .split(size);
        

            let btm_input_bar = Paragraph::new(input)
                .style(Style::default().fg(Color::LightCyan))
                .alignment(Alignment::Left)
                .block(
                    Block::default()
                        .borders(Borders::BOTTOM)
                        .style(Style::default().fg(App::border_colour(&view_state, &ViewState::Main)))
                        .title("CMD")
                        .border_type(BorderType::Plain),
                );


            let top_par = Paragraph::new(input)
                .style(Style::default().fg(Color::LightCyan))
                .alignment(Alignment::Left)
                .block(
                    Block::default()
                        .borders(Borders::BOTTOM)
                        .style(Style::default().fg(App::border_colour(&view_state, &ViewState::Workflows)))
                        .title(format!("{}", view_state))
                        .border_type(BorderType::Plain),
                );

            rect.render_widget(btm_input_bar, chunks[1]);
            rect.render_widget(top_par , chunks[0]);

        }).expect("hmm");
    }

    fn border_colour(view_state: &ViewState, this_view: &ViewState) -> Color {
        match view_state {
            ViewState::Main => 
                match this_view {
                    ViewState::Main => Color::LightCyan,
                    ViewState::Workflows => Color::LightGreen,
                },
            ViewState::Workflows => 
                match this_view {
                    ViewState::Main => Color::LightGreen,
                    ViewState::Workflows => Color::LightCyan,
                }
        }
    }

}




pub fn run() -> Result<(), Box<dyn std::error::Error>> {
    enable_raw_mode().expect("can run raw mode");

    let rx = spawn_input_thread();

    let stdout = io::stdout();
    let backend = CrosstermBackend::new(stdout);
    let mut terminal = Terminal::new(backend)?;
    let mut app = App::new(&mut terminal);
    app.terminal.clear()?;
    

    let mut input = String::new();
    loop {
        app.draw(&input);
        match app.handle_event(&rx)? {
            HandleEventMsg::Continue => continue,
            HandleEventMsg::Exit => break,
        }
    }

    Ok(())
}

/// Spawn an input thread and are polling for key events, 
/// returns a channel for receiving these inputs
fn spawn_input_thread() -> Receiver<Event<KeyEvent>> {
    let (tx, rx) = mpsc::channel();
    let tick_rate = Duration::from_millis(200);
    thread::spawn(move || {
        let mut last_tick = Instant::now();
        loop {
            let timeout = tick_rate
                .checked_sub(last_tick.elapsed())
                .unwrap_or_else(|| Duration::from_secs(0));

            if event::poll(timeout).expect("poll works") {
                if let CEvent::Key(key) = event::read().expect("can read events") {
                    tx.send(Event::Input(key)).expect("can send events");
                }
            }

            if last_tick.elapsed() >= tick_rate {
                if let Ok(_) = tx.send(Event::Tick) {
                    last_tick = Instant::now();
                }
            




            }
        }
    });

    rx
}