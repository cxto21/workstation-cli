use std::path::PathBuf;
use std::{fs, io};

pub fn get_socket_path() -> PathBuf {
    state_dir().join("daemon.sock")
}

pub fn get_log_path() -> PathBuf {
    state_dir().join("daemon.log")
}

pub fn get_client_log_path() -> PathBuf {
    state_dir().join("client.log")
}

pub fn get_lock_path() -> PathBuf {
    state_dir().join("daemon.lock")
}

pub fn get_pid_path() -> PathBuf {
    state_dir().join("daemon.pid")
}

pub fn get_state_file_path() -> PathBuf {
    config_dir().join("state.json")
}

pub fn get_config_file_path() -> PathBuf {
    config_dir().join("config.toml")
}

pub fn get_config_dir() -> PathBuf {
    config_dir()
}

pub fn migrate_legacy_user_data() {
    let new_cfg = config_dir();
    let new_state = state_dir();
    let legacy_cfg = legacy_config_dir();
    let legacy_state = legacy_state_dir();

    // Best-effort migration: never overwrite existing new files.
    let _ = migrate_file_if_missing(&legacy_cfg.join("config.toml"), &new_cfg.join("config.toml"));
    let _ = migrate_file_if_missing(&legacy_cfg.join("state.json"), &new_cfg.join("state.json"));
    let _ = migrate_file_if_missing(&legacy_cfg.join("theme.toml"), &new_cfg.join("theme.toml"));

    let _ = migrate_file_if_missing(&legacy_state.join("daemon.log"), &new_state.join("daemon.log"));
    let _ = migrate_file_if_missing(&legacy_state.join("client.log"), &new_state.join("client.log"));
}

fn state_dir() -> PathBuf {
    let dir = workstation_home().join("state");
    fs::create_dir_all(&dir).ok();
    dir
}

fn config_dir() -> PathBuf {
    let dir = workstation_home().join("config");
    fs::create_dir_all(&dir).ok();
    dir
}

fn workstation_home() -> PathBuf {
    if let Ok(root) = std::env::var("WORKSTATION_HOME") {
        return PathBuf::from(root);
    }
    let mut home = PathBuf::from(std::env::var("HOME").unwrap_or_else(|_| "/tmp".to_string()));
    home.push(".workstation");
    fs::create_dir_all(&home).ok();
    home
}

fn legacy_state_dir() -> PathBuf {
    if let Some(dirs) = directories::ProjectDirs::from("", "", "mato") {
        dirs.state_dir()
            .unwrap_or(dirs.config_dir())
            .to_path_buf()
    } else {
        PathBuf::from("/tmp/mato")
    }
}

fn legacy_config_dir() -> PathBuf {
    std::env::var("XDG_CONFIG_HOME")
        .map(PathBuf::from)
        .unwrap_or_else(|_| {
            let mut h = PathBuf::from(std::env::var("HOME").unwrap_or_default());
            h.push(".config");
            h
        })
        .join("mato")
}

fn migrate_file_if_missing(src: &PathBuf, dst: &PathBuf) -> io::Result<()> {
    if !src.exists() || dst.exists() {
        return Ok(());
    }
    if let Some(parent) = dst.parent() {
        fs::create_dir_all(parent)?;
    }
    let _ = fs::copy(src, dst)?;
    Ok(())
}
