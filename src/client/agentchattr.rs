use crate::client::app::ChatThread;
use serde::Deserialize;
use std::collections::HashMap;
use std::fs;
use std::path::{Path, PathBuf};
use std::time::{SystemTime, UNIX_EPOCH};

#[derive(Deserialize)]
struct AgentChattrMessage {
    channel: Option<String>,
    timestamp: Option<f64>,
}

pub fn load_chat_threads() -> Vec<ChatThread> {
    let Some(path) = resolve_log_path() else {
        return vec![];
    };
    let Ok(content) = fs::read_to_string(path) else {
        return vec![];
    };

    let mut channels: HashMap<String, f64> = HashMap::new();
    for line in content.lines() {
        let trimmed = line.trim();
        if trimmed.is_empty() {
            continue;
        }
        let Ok(msg) = serde_json::from_str::<AgentChattrMessage>(trimmed) else {
            continue;
        };
        let channel = msg.channel.unwrap_or_else(|| "general".to_string());
        let ts = msg.timestamp.unwrap_or(0.0);
        channels
            .entry(channel)
            .and_modify(|current| {
                if ts > *current {
                    *current = ts;
                }
            })
            .or_insert(ts);
    }

    let now = now_epoch_secs();
    let mut threads: Vec<(String, f64)> = channels.into_iter().collect();
    threads.sort_by(|a, b| b.1.partial_cmp(&a.1).unwrap_or(std::cmp::Ordering::Equal));

    threads
        .into_iter()
        .map(|(title, ts)| ChatThread {
            title,
            subtitle: relative_time(now, ts),
        })
        .collect()
}

fn resolve_log_path() -> Option<PathBuf> {
    // Explicit path override for integration/testing.
    if let Ok(explicit) = std::env::var("AGENTCHATTR_LOG_PATH") {
        let p = PathBuf::from(explicit);
        if p.exists() {
            return Some(p);
        }
    }

    // Data dir override used by agentchattr server config.
    if let Ok(data_dir) = std::env::var("AGENTCHATTR_DATA_DIR") {
        if let Some(p) = find_existing_log(Path::new(&data_dir)) {
            return Some(p);
        }
    }

    let mut candidates: Vec<PathBuf> = Vec::new();

    // Local repo execution defaults.
    candidates.push(PathBuf::from("./data"));

    if let Ok(home) = std::env::var("HOME") {
        candidates.push(PathBuf::from(&home).join(".local/share/agentchattr/data"));
        candidates.push(PathBuf::from(&home).join(".config/agentchattr/data"));
        candidates.push(PathBuf::from(&home).join(".agentchattr/data"));
    }

    for base in candidates {
        if let Some(p) = find_existing_log(&base) {
            return Some(p);
        }
    }

    None
}

fn find_existing_log(base: &Path) -> Option<PathBuf> {
    for name in ["agentchattr_log.jsonl", "room_log.jsonl"] {
        let p = base.join(name);
        if p.exists() {
            return Some(p);
        }
    }
    None
}

fn now_epoch_secs() -> f64 {
    match SystemTime::now().duration_since(UNIX_EPOCH) {
        Ok(d) => d.as_secs_f64(),
        Err(_) => 0.0,
    }
}

fn relative_time(now: f64, then: f64) -> String {
    if then <= 0.0 || now <= then {
        return "recently".to_string();
    }
    let delta = (now - then) as u64;
    if delta < 60 {
        format!("{}s ago", delta)
    } else if delta < 3600 {
        format!("{}m ago", delta / 60)
    } else if delta < 86_400 {
        format!("{}h ago", delta / 3600)
    } else {
        format!("{}d ago", delta / 86_400)
    }
}
