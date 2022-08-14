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
enum Views {
    Main,
    Workflows
}

impl Views {
    fn next(&self) -> Views {
        match self {
            Views::Main => Views::Workflows,
            Views::Workflows => Views::Main
        }
    }
}


impl Display for Views {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Views::Main => write!(f, "Main"),
            Views::Workflows => write!(f, "Workflows"),
        }
    }
}


struct StatefulList<T> {
    state: ListState,
    items: Vec<T>,
}

impl<T> StatefulList<T> {
    fn with_items(items: Vec<T>) -> Self {
        StatefulList {
            state: ListState::default(),
            items,
        }
    }

    fn next(&mut self) {
        let i = match self.state.selected() {
            Some(i) => {
                if i >= self.items.len() -1 {
                    0
                } else {
                    i + 1
                }
            },
            None => 0,
        };
        self.state.select(Some(i));


    }

    fn previous(&mut self) {
        let i = match self.state.selected() {
            Some(i) => {
                if i == 0 {
                    self.items.len() - 1
                } else {
                    i - 1
                }
            }
            None => 0,
        };
        self.state.select(Some(i));
    }

    fn unselect(&mut self) {
        self.state.select(None);
    }


}

struct App<'a> {
    view_state: Views,
    terminal: &'a mut Terminal<CrosstermBackend<io::Stdout>>,
    workflows_items: StatefulList<(&'a str, &'a str)>,
}

impl App<'_> {
    fn new(terminal: &mut Terminal<CrosstermBackend<io::Stdout>>) -> App {
        App {
            view_state: Views::Main,
            terminal: terminal,
            workflows_items: StatefulList::with_items(vec![("Workflow 1", "id:1234"), ("Workflow 2", "id:1552")]),
        
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
                .constraints(
                    [
                        Constraint::Min(2),
                        Constraint::Length(3),
                    ]
                    .as_ref(),
                )
                .split(size);
        

            


                let items: Vec<ListItem> = self.workflows_items.items.iter().map(|(name, id)| {
                    ListItem::new(Spans::from(vec![Span::raw(format!("{} - {}", name, id))]))
                }).collect();

                let items = List::new(items)
                    .block(Block::default().borders(Borders::ALL).title("List"))
                    .highlight_style(
                        Style::default()
                            .bg(Color::LightGreen)
                            .add_modifier(Modifier::BOLD),
                    )
                    .highlight_symbol(">> ");

            // We can now render the item list
            rect.render_stateful_widget(items, chunks[0], &mut self.workflows_items.state);


            let btm_input_bar = Paragraph::new(input)
                .style(Style::default().fg(Color::LightCyan))
                .alignment(Alignment::Left)
                .block(
                    Block::default()
                        .borders(Borders::BOTTOM)
                        .style(Style::default().fg(App::border_colour(&view_state, &Views::Main)))
                        .title("CMD")
                        .border_type(BorderType::Plain),
                );


            let top_par = Paragraph::new(input)
                .style(Style::default().fg(Color::LightCyan))
                .alignment(Alignment::Left)
                .block(
                    Block::default()
                        .borders(Borders::BOTTOM)
                        .style(Style::default().fg(App::border_colour(&view_state, &Views::Workflows)))
                        .title(format!("{}", view_state))
                        .border_type(BorderType::Plain),
                );

            rect.render_widget(btm_input_bar, chunks[1]);
            rect.render_widget(top_par , chunks[0]);

        }).expect("hmm");
    }

    fn border_colour(view_state: &Views, this_view: &Views) -> Color {
        match view_state {
            Views::Main => 
                match this_view {
                    Views::Main => Color::LightRed,
                    Views::Workflows => Color::Gray,
                },
            Views::Workflows => 
                match this_view {
                    Views::Main => Color::Gray,
                    Views::Workflows => Color::LightRed,
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