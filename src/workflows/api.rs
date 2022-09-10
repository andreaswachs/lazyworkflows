use reqwest::{Client, Request, RequestBuilder};
use crate::config;
use super::dsl::{self, ReponseSerializable};

#[derive(Clone, Debug, Copy)]
enum Action {
    List,
    Get,
    Dispatch,
    Enable,
    Disable,
}

#[derive(Debug, Clone)]
struct APIRequest {
    action: Option<Action>,
    owner: Option<String>,
    repo: Option<String>,
    id: Option<String>,
    token: Option<String>,
    client: Client,
}

impl APIRequest {
    fn new() -> APIRequest {
        APIRequest {
            action: None,
            owner: None,
            repo: None, 
            id: None,
            token: None, 
            client: Client::new(),
        }
    }

    fn set_action(&mut self, action: &Action) -> &mut Self {
        self.action = Some(action.clone());
        self
    }

    fn set_owner(&mut self, owner: &String) -> &mut Self {
        self.owner = Some(owner.clone());
        self
    }

    fn set_repo(&mut self, reponame: &String) -> &mut Self {
        self.repo = Some(reponame.clone());
        self
    }

    fn set_id(&mut self, id: &String) -> &mut Self {
        self.id = Some(id.clone());
        self
    }

    fn set_token(&mut self, token: &String) -> &mut Self {
        self.token = Some(token.clone());
        self
    }

    fn use_repo(&mut self, repo: &config::manager::Repo) -> &mut Self {
        self.set_owner(&repo.owner)
            .set_repo(&repo.repo)
            .set_token(&repo.token);
        self
    }

    fn generate_url(&self) -> String {
        match self.action {
            Some(action) => 
                match action {
                    Action::List => format!("https://api.github.com/repos/{}/{}/actions/workflows", self.owner.as_ref().unwrap(), self.repo.as_ref().unwrap()),
                    Action::Get => format!("https://api.github.com/repos/{}/{}/actions/workflows/{}", self.owner.as_ref().unwrap(), self.repo.as_ref().unwrap(), self.id.as_ref().unwrap()),
                    Action::Dispatch => format!("https://api.github.com/repos/{}/{}/actions/workflows/{}/dispatches", self.owner.as_ref().unwrap(), self.repo.as_ref().unwrap(), self.id.as_ref().unwrap()),
                    Action::Enable => format!("https://api.github.com/repos/{}/{}/actions/workflows/{}/enable", self.owner.as_ref().unwrap(), self.repo.as_ref().unwrap(), self.id.as_ref().unwrap()),
                    Action::Disable => format!("https://api.github.com/repos/{}/{}/actions/workflows/{}/disable", self.owner.as_ref().unwrap(), self.repo.as_ref().unwrap(), self.id.as_ref().unwrap()),
                }
            None => panic!("No action set"),
        }
    }

    fn apply_http_method(&self) -> RequestBuilder {
        let url = self.generate_url();
        match self.action {
            Some(action) => match action {
                Action::List => self.client.get(url),
                Action::Get => self.client.get(url),
                Action::Dispatch => self.client.post(url),
                Action::Enable => self.client.put(url),
                Action::Disable => self.client.put(url),
            },
            None => panic!("No action set"),
        }
    }

    fn build(&self) -> Request {
        // TODO: we could build this nicer by moving the header building into a separate function,
        // but id prefer it to be located within this impl, but that seems to be a little difficult
        let mut request = 
            self.apply_http_method()
            .header("Accept", "application/vnd.github.v3+json")
            .header("Authorization", format!("token {}", &self.token.as_ref().unwrap()))
            .header("User-Agent", "lazyworkflows");

        request.build().unwrap()
    }

    async fn send(&self) -> Result<String, reqwest::Error> {
        let request = self.build();
        let response = self.client.execute(request).await.unwrap().text().await.unwrap();

        Ok(response)
    }
}

pub async fn list(cfg: &config::manager::Repo) -> dsl::ListResponse{
    let response = APIRequest::new()
        .set_action(&Action::List)
        .use_repo(&cfg)
        .send()
        .await
        .unwrap();

    dsl::ListResponse::serialize_from(&response)
}

pub async fn get(cfg: &config::manager::Repo, id: &String) -> dsl::GetResponse {
    let response = 
    APIRequest::new()
        .set_action(&Action::Get)
        .use_repo(&cfg)
        .set_id(id)
        .send()
        .await
        .unwrap();

    dsl::GetResponse::serialize_from(&response)
}

pub async fn dispatch(cfg: &config::manager::Repo, id: &String) -> dsl::DispatchResponse {
    // TODO: Add support for input
    let response = APIRequest::new()
        .set_action(&Action::Dispatch)
        .use_repo(&cfg)
        .set_id(id)
        .send()
        .await
        .unwrap();

    dsl::DispatchResponse::serialize_from(&response)
}

pub async fn enable(cfg: &config::manager::Repo, id: &String) -> dsl::EnableResponse {
    let response = APIRequest::new()
        .set_action(&Action::Enable)
        .use_repo(&cfg)
        .set_id(id)
        .send()
        .await
        .unwrap();

    dsl::EnableResponse::serialize_from(&response)
}

pub async fn disable(cfg: &config::manager::Repo, id: &String) -> dsl::DisableResponse {
    let response = APIRequest::new()
        .set_action(&Action::Disable)
        .use_repo(&cfg)
        .set_id(id)
        .send()
        .await
        .unwrap();
    
    dsl::DisableResponse::serialize_from(&response)
}
