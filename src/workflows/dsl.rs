use chrono::{DateTime, Utc, FixedOffset, format::Fixed};
use serde::{Deserialize, Serialize, Deserializer};

//
// Get reponse
//

#[derive(Default, Serialize, Deserialize, Debug)]
pub struct WorkflowReponse {
    id: i32,
    node_id: String,
    name: String,
    path: String,
    state: String,
    created_at: String,
    updated_at: String,
    url: String,
    html_url: String,
    badge_url: String
}

#[derive(Serialize, Deserialize, Debug)]
pub struct GetResponse {
    workflow: WorkflowReponse,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct ListResponse {
    total_count: i32,
    workflows: Vec<WorkflowReponse>,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct EnableResponse {
    status: i32
}

#[derive(Serialize, Deserialize, Debug)]
pub struct DisableResponse {
    status: i32
}

#[derive(Serialize, Deserialize, Debug)]
pub struct DispatchResponse {
    status: i32
}

fn serialize<'a, T: Deserialize<'a>>(input: &'a String) -> T {
    let response: T = serde_json::from_str(input).unwrap();
    response
}

pub fn serialize_get(input: &String) -> GetResponse {
    GetResponse {
        workflow: serialize::<WorkflowReponse>(input),
    }
}

pub fn serialize_list(input: &String) -> ListResponse {
    serialize::<ListResponse>(input)
}

pub fn serialize_enable(input: &String) -> EnableResponse {
    serialize::<EnableResponse>(input)
}

pub fn serialize_disable(input: &String) -> DisableResponse {
    serialize::<DisableResponse>(input)
}

pub fn serialize_dispatch(input: &String) -> DispatchResponse {
    serialize::<DispatchResponse>(input)
}



#[cfg(test)]
mod tests {
    use super::*;

    fn workflow1<'a>() -> &'a str {
        r#"{"id":161335,"node_id":"MDg6V29ya2Zsb3cxNjEzMzU=","name":"CI","path":".github/workflows/blank.yaml","state":"active","created_at":"2020-01-08T23:48:37.000-08:00","updated_at":"2020-01-08T23:50:21.000-08:00","url":"https://api.github.com/repos/octo-org/octo-repo/actions/workflows/161335","html_url":"https://github.com/octo-org/octo-repo/blob/master/.github/workflows/161335","badge_url":"https://github.com/octo-org/octo-repo/workflows/CI/badge.svg"}"#
    }

    fn workflow2<'a>() -> &'a str {
        r#"{"id":20,"node_id":"MDg6V29ya2Zsb3cxNjEzMzU=","name":"CD","path":".github/workflows/other.yaml","state":"disabled","created_at":"2020-01-08T23:48:37.000-08:00","updated_at":"2020-01-08T23:50:21.000-08:00","url":"https://api.github.com/repos/octo-org/octo-repo/actions/workflows/161335","html_url":"https://github.com/octo-org/octo-repo/blob/master/.github/workflows/161335","badge_url":"https://github.com/octo-org/octo-repo/workflows/CI/badge.svg"}"#
    }

    #[test]
    fn serialize_enable_response() {
        let input = r#"{"status": 200}"#;
        let response = serialize_enable(&input.to_string());
        assert_eq!(response.status, 200);
    }

    #[test]
    fn serialize_disable_response() {
        let input = r#"{"status": 200}"#;
        let response = serialize_disable(&input.to_string());
        assert_eq!(response.status, 200);
    }

    #[test]
    fn serialize_dispatch_response() {
        let input = r#"{"status": 200}"#;
        let response = serialize_dispatch(&input.to_string());
        assert_eq!(response.status, 200);
    }

    #[test]
    fn serialize_get_response() {
        // The input is the exampple response from https://docs.github.com/en/rest/actions/workflows#get-a-workflow
        let input = workflow1();

        let response = serialize_get(&input.to_string());
        let workflow = response.workflow;

        assert_eq!(workflow.id, 161335);
        assert_eq!(&workflow.name, "CI");
        assert_eq!(&workflow.path, ".github/workflows/blank.yaml");
        assert_eq!(&workflow.state, "active");
        assert_eq!(&workflow.created_at, "2020-01-08T23:48:37.000-08:00");
        assert_eq!(&workflow.updated_at, "2020-01-08T23:50:21.000-08:00");
        assert_eq!(&workflow.url, "https://api.github.com/repos/octo-org/octo-repo/actions/workflows/161335");
        assert_eq!(&workflow.html_url, "https://github.com/octo-org/octo-repo/blob/master/.github/workflows/161335");
        assert_eq!(&workflow.badge_url, "https://github.com/octo-org/octo-repo/workflows/CI/badge.svg");
    }


    #[test]
    fn serialize_list_response() {
        let mut input = String::new();
        input.push_str(r#"{"total_count": 2, "workflows": ["#);
        input.push_str(workflow1());
        input.push_str(",");
        input.push_str(workflow2());
        input.push_str("]}");

        let response = serialize_list(&input.to_string());
        let workflows = response.workflows;

        assert_eq!(workflows.len(), 2);
        assert_eq!(workflows[0].id, 161335);
        assert_eq!(&workflows[0].name, "CI");
        assert_eq!(&workflows[0].path, ".github/workflows/blank.yaml");
        assert_eq!(&workflows[0].state, "active");
        assert_eq!(&workflows[0].created_at, "2020-01-08T23:48:37.000-08:00");
        assert_eq!(&workflows[0].updated_at, "2020-01-08T23:50:21.000-08:00");
        assert_eq!(&workflows[0].url, "https://api.github.com/repos/octo-org/octo-repo/actions/workflows/161335");
        assert_eq!(&workflows[0].html_url, "https://github.com/octo-org/octo-repo/blob/master/.github/workflows/161335");
        assert_eq!(&workflows[0].badge_url, "https://github.com/octo-org/octo-repo/workflows/CI/badge.svg");
        assert_eq!(workflows[1].id, 20);
        assert_eq!(&workflows[1].name, "CD");
        assert_eq!(&workflows[1].path, ".github/workflows/other.yaml");
        assert_eq!(&workflows[1].state, "disabled");
        assert_eq!(&workflows[1].created_at, "2020-01-08T23:48:37.000-08:00");
        assert_eq!(&workflows[1].updated_at, "2020-01-08T23:50:21.000-08:00");
        assert_eq!(&workflows[1].url, "https://api.github.com/repos/octo-org/octo-repo/actions/workflows/161335");
        assert_eq!(&workflows[1].html_url, "https://github.com/octo-org/octo-repo/blob/master/.github/workflows/161335");
        assert_eq!(&workflows[1].badge_url, "https://github.com/octo-org/octo-repo/workflows/CI/badge.svg");
    }
}