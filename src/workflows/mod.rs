
use super::config;

mod api;
mod dsl;

pub struct Workflow {
    name: Option<String>,
    path: Option<String>,
    id: Option<String>,
    active: Option<String> // TODO: make this an enum with enum -> string conversion
}

impl Workflow {
    pub fn new() -> Workflow {
        Workflow {
            name: None,
            path: None,
            id: None,
            active: None
        }
    }

    pub fn set_name<'a>(&'a mut self, name: String) -> &'a mut Workflow {
        self.name = Some(name);
        self
    }

    pub fn set_path<'a>(&'a mut self, path: String) -> &'a mut Workflow {
        self.path = Some(path);
        self
    }

    pub fn set_id<'a>(&'a mut self, id: String) -> &'a mut Workflow {
        self.id = Some(id);
        self
    }

    pub fn set_active<'a>(&'a mut self, active: String) -> &'a mut Workflow {
        self.active = Some(active);
        self
    }
}

pub struct Repo {
    name: String,
    workflows: Vec<Workflow>,
}


pub struct Owner {
    name: String,
    repos: Vec<Repo>,
}

pub struct Workflows {
    owners: Vec<Owner>,
}

/*
   Public API functions regarding ALL workflows
*/

pub async fn load(config: &config::Config) -> Workflows {
    let mut requests = Vec::new();
    for repo in config.repos.iter() {
        requests.push(api::list(&repo));
    }

    let finished_requests = futures::future::join_all(requests).await;

    for request in finished_requests {
        println!("{:?}", request);
    }

    Workflows {
        owners: Vec::new(),
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn new_should_return_workflow_with_no_fields_set() {
        let workflow = Workflow::new();
        assert_eq!(workflow.name, None);
        assert_eq!(workflow.path, None);
        assert_eq!(workflow.id, None);
        assert_eq!(workflow.active, None);
    }

    #[test]
    fn set_name_should_set_name_field() {
        let mut workflow = Workflow::new();
        workflow.set_name("test".to_string());
        assert_eq!(workflow.name, Some("test".to_string()));
    }

    #[test]
    fn set_path_should_set_path_field() {
        let mut workflow = Workflow::new();
        workflow.set_path("test".to_string());
        assert_eq!(workflow.path, Some("test".to_string()));
    }

    #[test]
    fn set_id_should_set_id_field() {
        let mut workflow = Workflow::new();
        workflow.set_id("test".to_string());
        assert_eq!(workflow.id, Some("test".to_string()));
    }

    #[test]
    fn set_active_should_set_active_field() {
        let mut workflow = Workflow::new();
        workflow.set_active("test".to_string());
        assert_eq!(workflow.active, Some("test".to_string()));
    }


    #[test]
    fn create_full_workflow_should_contain_all_fields() {
        let mut workflow = Workflow::new();
        workflow.set_name("name".to_string());
        workflow.set_path("path".to_string());
        workflow.set_id("id".to_string());
        workflow.set_active("active".to_string());
        assert_eq!(workflow.name, Some("name".to_string()));
        assert_eq!(workflow.path, Some("path".to_string()));
        assert_eq!(workflow.id, Some("id".to_string()));
        assert_eq!(workflow.active, Some("active".to_string()));
    }
}