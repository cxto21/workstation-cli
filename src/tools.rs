use crate::error::{MatoError, Result};
use std::io::{self, Write};
use std::process::Command;

const GENTLE_AI_INSTALL_HINT: &str =
    "curl -fsSL https://raw.githubusercontent.com/Gentleman-Programming/gentle-ai/main/scripts/install.sh | bash";

pub fn handle_tools_command(args: &[String]) -> Result<()> {
    match args.first().map(|s| s.as_str()) {
        None | Some("help") | Some("--help") | Some("-h") => {
            print_tools_help();
            Ok(())
        }
        Some("status") => cmd_tools_status(),
        Some("init") => cmd_tools_init(),
        Some("sync") => cmd_tools_sync(),
        Some(other) => {
            eprintln!("Unknown tools command: {}", other);
            eprintln!();
            print_tools_help();
            Err(MatoError::StateSaveFailed(
                "Invalid tools subcommand".to_string(),
            ))
        }
    }
}

fn print_tools_help() {
    println!(
        "workstation-cli tools <command>\n\
         \n\
         Commands:\n\
           status   Show gentle-ai availability and version\n\
           init     Guided bootstrap (dry-run + optional apply)\n\
           sync     Run gentle-ai sync\n\
           help     Show this help\n"
    );
}

fn cmd_tools_status() -> Result<()> {
    let installed = gentle_ai_installed();
    if !installed {
        println!("gentle-ai: not found in PATH");
        println!("Install hint:");
        println!("  {}", GENTLE_AI_INSTALL_HINT);
        return Ok(());
    }

    let version = gentle_ai_version().unwrap_or_else(|| "unknown".to_string());
    println!("gentle-ai: available");
    println!("version: {}", version);
    println!("ready: run 'workstation-cli tools init' for guided bootstrap");
    Ok(())
}

fn cmd_tools_init() -> Result<()> {
    if !gentle_ai_installed() {
        println!("gentle-ai is not installed.");
        println!("Install hint:");
        println!("  {}", GENTLE_AI_INSTALL_HINT);
        return Err(MatoError::StateSaveFailed(
            "gentle-ai not available in PATH".to_string(),
        ));
    }

    println!("Workstation + gentle-ai bootstrap");
    println!("This guided setup will run a dry-run first, then ask before applying.");
    println!();

    let agents = prompt_with_default("Agents (comma-separated)", "claude-code,opencode")?;
    let preset = prompt_with_default("Preset", "ecosystem-only")?;
    let persona = prompt_with_default("Persona", "neutral")?;

    println!();
    println!("Step 1/2: dry-run");
    run_gentle_ai_install(&agents, &preset, &persona, true)?;

    println!();
    if !prompt_yes_no("Dry-run passed. Apply now?", true)? {
        println!("Bootstrap stopped after dry-run.");
        println!("When ready, run: workstation-cli tools init");
        return Ok(());
    }

    println!("Step 2/2: apply");
    run_gentle_ai_install(&agents, &preset, &persona, false)?;

    println!();
    println!("Bootstrap completed.");
    println!("Next: run 'workstation-cli tools sync' to re-apply ecosystem updates when needed.");
    Ok(())
}

fn cmd_tools_sync() -> Result<()> {
    if !gentle_ai_installed() {
        println!("gentle-ai is not installed.");
        println!("Install hint:");
        println!("  {}", GENTLE_AI_INSTALL_HINT);
        return Err(MatoError::StateSaveFailed(
            "gentle-ai not available in PATH".to_string(),
        ));
    }

    let status = Command::new("gentle-ai")
        .arg("sync")
        .status()
        .map_err(MatoError::from)?;

    if !status.success() {
        return Err(MatoError::StateSaveFailed(
            "gentle-ai sync failed".to_string(),
        ));
    }

    println!("gentle-ai sync completed.");
    Ok(())
}

fn run_gentle_ai_install(agents: &str, preset: &str, persona: &str, dry_run: bool) -> Result<()> {
    let mut cmd = Command::new("gentle-ai");
    cmd.arg("install")
        .arg("--agent")
        .arg(agents)
        .arg("--preset")
        .arg(preset)
        .arg("--persona")
        .arg(persona);

    if dry_run {
        cmd.arg("--dry-run");
    }

    let status = cmd.status().map_err(MatoError::from)?;
    if !status.success() {
        let mode = if dry_run { "dry-run" } else { "apply" };
        return Err(MatoError::StateSaveFailed(format!(
            "gentle-ai install failed during {}",
            mode
        )));
    }

    Ok(())
}

fn gentle_ai_installed() -> bool {
    Command::new("gentle-ai")
        .arg("version")
        .output()
        .map(|o| o.status.success())
        .unwrap_or(false)
}

fn gentle_ai_version() -> Option<String> {
    let output = Command::new("gentle-ai").arg("version").output().ok()?;
    if !output.status.success() {
        return None;
    }

    let stdout = String::from_utf8_lossy(&output.stdout).trim().to_string();
    if stdout.is_empty() {
        None
    } else {
        Some(stdout)
    }
}

fn prompt_with_default(label: &str, default_value: &str) -> Result<String> {
    print!("{} [{}]: ", label, default_value);
    io::stdout().flush().map_err(MatoError::from)?;

    let mut input = String::new();
    io::stdin().read_line(&mut input).map_err(MatoError::from)?;
    let value = input.trim();
    if value.is_empty() {
        Ok(default_value.to_string())
    } else {
        Ok(value.to_string())
    }
}

fn prompt_yes_no(label: &str, default_yes: bool) -> Result<bool> {
    let suffix = if default_yes { "[Y/n]" } else { "[y/N]" };
    print!("{} {}: ", label, suffix);
    io::stdout().flush().map_err(MatoError::from)?;

    let mut input = String::new();
    io::stdin().read_line(&mut input).map_err(MatoError::from)?;
    let answer = input.trim().to_ascii_lowercase();
    if answer.is_empty() {
        return Ok(default_yes);
    }
    Ok(answer == "y" || answer == "yes")
}