<!doctype html>
<html lang="ru">
<head>
    {{template "head.tpl"}}
    <title>{{.Post.Title}}</title>
</head>
<body class="uk-background-muted">
<div class="uk-margin-left uk-margin-right uk-margin-top">
    <div class="uk-card uk-card-default uk-card-hover uk-width-1-1 uk-margin-bottom">
        <div class="uk-card-header">
            <div class="uk-flex uk-flex-right">
                {{template "tools.tpl" .Post}}
            </div>
            <div class="uk-grid-small uk-flex-middle" uk-grid>
                <div class="uk-width-auto">
                    <img class="uk-border-circle" width="40" height="40" src="../static/img/avatar.png">
                </div>
                <div class="uk-width-expand">
                    <h2 class="uk-card-title ">{{.Post.Title}}</h2>
                    <p class="uk-text-meta">{{.Post.Date2Norm}}</p>
                </div>
            </div>
        </div>
        <div class="uk-card-body">
            <p>{{.Post.Summary}}</p>
        </div>
        <div class="uk-card-footer">
            <p>{{.Post.Body}}</p>
        </div>
        <div class="uk-card-footer uk-flex uk-flex-between">
            <a class="uk-link-heading uk-text-primary uk-button-text" href="/"><< Вернуться</a>
        </div>
    </div>
</div>
{{template "footer.tpl"}}
</body>
</html>
