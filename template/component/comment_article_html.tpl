{{define "comment_article_html"}}

{{if eq .config.CommentProvider "duoshuo"}}
{{template "comment_article_html_duoshuo"}}
{{else if eq .config.CommentProvider "livere"}}
{{template "comment_article_html_livere"}}
{{else if eq .config.CommentProvider "163"}}
{{template "comment_article_html_163"}}
{{else}}
Comment disabled ({{.config.CommentProvider}})
{{end}}

{{end}}