use serde::{Deserialize, Serialize};
use serde_json::Result;

#[derive(Serialize, Deserialize)]
pub struct JsonAction {
    pub action: String,
    pub path: String
}

#[derive(Serialize, Deserialize)]
pub struct JsonFile {
    pub status: String
}

#[derive(Serialize, Deserialize)]
pub struct JsonIssue {
    pub severity: String,
    pub message: String,
    pub location: String
}

#[derive(Serialize, Deserialize)]
pub struct JsonLog {
    pub actions: Vec<JsonAction>,
    pub issues: Vec<JsonIssue>
}

#[derive(Serialize, Deserialize)]
pub struct JsonLogger(Vec<JsonLog>);

