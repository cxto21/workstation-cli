use crate::client::app::{App, Focus, WorkspaceView};
use crate::theme::ThemeColors;
use ratatui::{
    layout::{Constraint, Direction, Layout, Rect},
    style::{Modifier, Style},
    text::{Line, Span},
    widgets::{Block, Borders, List, ListItem, ListState, Paragraph, Wrap},
    Frame,
};

pub(super) fn draw_workspace(f: &mut Frame, app: &mut App, area: Rect, t: &ThemeColors) {
    app.kb_list_area = Rect::default();
    match app.workspace_view {
        WorkspaceView::Chat => draw_chat_layout(f, app, area, t),
        WorkspaceView::Alerts => draw_simple_layout(f, app, area, t, "Alerts"),
        WorkspaceView::Logs => draw_simple_layout(f, app, area, t, "Logs"),
        WorkspaceView::Docs => draw_docs_layout(f, app, area, t),
        WorkspaceView::Terminal => {}
    }
}

fn draw_chat_layout(f: &mut Frame, app: &mut App, area: Rect, t: &ThemeColors) {
    let active = app.focus == Focus::Content;
    let left_rows = Layout::default()
        .direction(Direction::Vertical)
        .constraints([Constraint::Length(4), Constraint::Min(0)])
        .split(area);

    let new_chat_style = if t.follow_terminal {
        Style::default().add_modifier(Modifier::DIM)
    } else {
        Style::default().fg(t.fg_dim()).bg(t.surface()).add_modifier(Modifier::DIM)
    };
    f.render_widget(
        Paragraph::new(Line::from(vec![Span::styled(
            " + NEW CHAT (IN DEVELOPMENT) ",
            new_chat_style,
        )]))
            .wrap(Wrap { trim: false })
            .block(
                Block::default()
                    .borders(Borders::ALL)
                    .title(Span::styled(
                        " CHAT HISTORY ",
                        super::title_style(t, active),
                    ))
                    .border_style(super::border_style(t, active))
                    .style(Style::default().bg(t.surface())),
            ),
        left_rows[0],
    );

    let items: Vec<ListItem> = app
        .chat_threads
        .iter()
        .map(|thread| {
            ListItem::new(Line::from(vec![
                Span::styled(thread.title.clone(), Style::default().fg(t.fg())),
                Span::raw("\n"),
                Span::styled(thread.subtitle.clone(), Style::default().fg(t.fg_dim())),
            ]))
        })
        .collect();

    if app.chat_threads.is_empty() {
        f.render_widget(
            Paragraph::new("No chats yet. Chat integration is in development.")
                .style(Style::default().fg(t.fg_dim()))
                .block(
                    Block::default()
                        .borders(Borders::ALL)
                        .title(Span::styled(" CHATS ", super::title_style(t, active)))
                        .border_style(super::border_style(t, active))
                        .style(Style::default().bg(t.surface())),
                ),
            left_rows[1],
        );
    } else {
        let mut list_state = ListState::default();
        list_state.select(Some(app.chat_selected.min(app.chat_threads.len() - 1)));
        f.render_stateful_widget(
            List::new(items)
                .block(
                    Block::default()
                        .borders(Borders::ALL)
                        .title(Span::styled(" CHATS ", super::title_style(t, active)))
                        .border_style(super::border_style(t, active))
                        .style(Style::default().bg(t.surface())),
                )
                .highlight_style(Style::default().bg(t.sel_bg()).add_modifier(Modifier::BOLD)),
            left_rows[1],
            &mut list_state,
        );
    }
}

fn draw_simple_layout(f: &mut Frame, app: &App, area: Rect, t: &ThemeColors, title: &str) {
    let active = app.focus == Focus::Content;
    let rows = Layout::default()
        .direction(Direction::Vertical)
        .constraints([Constraint::Length(3), Constraint::Min(0)])
        .split(area);

    f.render_widget(
        Block::default()
            .borders(Borders::ALL)
            .title(Span::styled(
                format!(" {} ", title.to_uppercase()),
                super::title_style(t, active),
            ))
            .border_style(super::border_style(t, active))
            .style(Style::default().bg(t.surface())),
        rows[0],
    );

    f.render_widget(
        Block::default()
            .borders(Borders::ALL)
            .title(Span::styled(" BODY ", super::title_style(t, active)))
            .border_style(super::border_style(t, active))
            .style(Style::default().bg(t.bg())),
        rows[1],
    );
}

fn draw_docs_layout(f: &mut Frame, app: &mut App, area: Rect, t: &ThemeColors) {
    let active = app.focus == Focus::Content;
    let cols = Layout::default()
        .direction(Direction::Horizontal)
        .constraints([Constraint::Length(34), Constraint::Min(0)])
        .split(area);

    let left_rows = Layout::default()
        .direction(Direction::Vertical)
        .constraints([Constraint::Length(3), Constraint::Min(0)])
        .split(cols[0]);

    f.render_widget(
        Paragraph::new("Knowledge Base (Markdown)")
            .style(Style::default().fg(t.fg_dim()))
            .block(
                Block::default()
                    .borders(Borders::ALL)
                    .title(Span::styled(" KB DOCS ", super::title_style(t, active)))
                    .border_style(super::border_style(t, active))
                    .style(Style::default().bg(t.surface())),
            ),
        left_rows[0],
    );

    app.kb_list_area = left_rows[1];
    if app.kb_docs.is_empty() {
        f.render_widget(
            Paragraph::new("No markdown files found under docs/")
                .style(Style::default().fg(t.fg_dim()))
                .wrap(Wrap { trim: false })
                .block(
                    Block::default()
                        .borders(Borders::ALL)
                        .title(Span::styled(" FILES ", super::title_style(t, active)))
                        .border_style(super::border_style(t, active))
                        .style(Style::default().bg(t.surface())),
                ),
            left_rows[1],
        );
    } else {
        let items: Vec<ListItem> = app
            .kb_docs
            .iter()
            .map(|doc| ListItem::new(doc.title.clone()))
            .collect();

        let mut list_state = ListState::default();
        let selected = app.kb_selected.min(app.kb_docs.len() - 1);
        app.kb_selected = selected;
        list_state.select(Some(selected));

        f.render_stateful_widget(
            List::new(items)
                .block(
                    Block::default()
                        .borders(Borders::ALL)
                        .title(Span::styled(" FILES ", super::title_style(t, active)))
                        .border_style(super::border_style(t, active))
                        .style(Style::default().bg(t.surface())),
                )
                .highlight_style(Style::default().bg(t.sel_bg()).add_modifier(Modifier::BOLD)),
            left_rows[1],
            &mut list_state,
        );
    }

    let right_rows = Layout::default()
        .direction(Direction::Vertical)
        .constraints([Constraint::Length(3), Constraint::Min(0)])
        .split(cols[1]);

    let current_path = app
        .kb_docs
        .get(app.kb_selected)
        .map(|d| d.path.clone())
        .unwrap_or_else(|| "docs/".to_string());
    f.render_widget(
        Paragraph::new(current_path)
            .style(Style::default().fg(t.fg_dim()))
            .block(
                Block::default()
                    .borders(Borders::ALL)
                    .title(Span::styled(" PREVIEW ", super::title_style(t, active)))
                    .border_style(super::border_style(t, active))
                    .style(Style::default().bg(t.surface())),
            ),
        right_rows[0],
    );

    let content = app
        .kb_docs
        .get(app.kb_selected)
        .map(|d| d.content.clone())
        .unwrap_or_else(|| "No content available.".to_string());
    f.render_widget(
        Paragraph::new(content)
            .style(Style::default().fg(t.fg()))
            .wrap(Wrap { trim: false })
            .block(
                Block::default()
                    .borders(Borders::ALL)
                    .title(Span::styled(" CONTENT ", super::title_style(t, active)))
                    .border_style(super::border_style(t, active))
                    .style(Style::default().bg(t.bg())),
            ),
        right_rows[1],
    );
}
