use crate::client::app::KbDoc;
use std::fs;
use std::path::{Path, PathBuf};

pub fn load_kb_docs() -> Vec<KbDoc> {
    let root = PathBuf::from("docs");
    if !root.exists() {
        return vec![];
    }

    let mut markdown_files: Vec<PathBuf> = vec![];
    collect_markdown_files(&root, &mut markdown_files);
    markdown_files.sort();

    markdown_files
        .into_iter()
        .filter_map(|path| to_kb_doc(&path, &root))
        .collect()
}

fn collect_markdown_files(dir: &Path, out: &mut Vec<PathBuf>) {
    let Ok(entries) = fs::read_dir(dir) else {
        return;
    };
    for entry in entries.flatten() {
        let path = entry.path();
        if path.is_dir() {
            collect_markdown_files(&path, out);
            continue;
        }
        if path
            .extension()
            .and_then(|e| e.to_str())
            .map(|e| e.eq_ignore_ascii_case("md"))
            .unwrap_or(false)
        {
            out.push(path);
        }
    }
}

fn to_kb_doc(path: &Path, root: &Path) -> Option<KbDoc> {
    let content = fs::read_to_string(path).ok()?;
    let rel = path.strip_prefix(root).ok()?.to_string_lossy().to_string();
    let title = rel.clone();
    let content = if content.len() > 100_000 {
        content[..100_000].to_string()
    } else {
        content
    };
    Some(KbDoc {
        title,
        path: rel,
        content,
    })
}
