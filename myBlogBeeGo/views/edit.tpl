<!doctype html>
<html lang="ru">
<head>
    {{template "head.tpl"}}
    <title>{{.Post.Title}}</title>
</head>
<body class="uk-background-muted">
<form onsubmit="updatePost('{{.Post.ID}}')">
    <fieldset id="editPost" class="uk-fieldset">
        <div class="uk-margin-left uk-margin-right uk-margin-top" uk-grid>
            <div class="uk-card uk-card-default uk-card-hover uk-width-1-1">
                <div class="uk-card-header">
                    <div class="uk-grid-small uk-flex-middle" uk-grid>
                        <div class="uk-width-auto">
                            <img class="uk-border-circle" width="40" height="40"
                                 src="../../static/img/avatar.png">
                        </div>
                        <div class="uk-width-expand">
                            <h3><input required type="text" class="uk-input" name="title" id="title"
                                       value="{{.Post.Title}}"></h3>
                            <input required class="uk-text-meta" type="text" class="uk-input" name="date" id="date"
                                   value="{{.Post.Date2Norm}}" readonly>
                        </div>
                    </div>
                </div>
                <div class="uk-card-body">
                        <textarea required class="uk-textarea" rows="5" name="summary"
                                  id="summary">{{.Post.Summary}}</textarea>
                </div>
                <div class="uk-card-footer">
                    <textarea required class="uk-textarea" rows="10" name="body" id="body">{{.Post.Body}}</textarea>
                </div>
                <div class="uk-card-footer">
                    <a class="uk-link-heading uk-text-primary uk-button-text" href="/"><< Вернуться</a>
                </div>
            </div>
        </div>
    </fieldset>
    <div class="uk-flex uk-flex-right uk-margin-small-top uk-margin-right uk-margin-small-bottom">
        <button type="submit" class="uk-button uk-button-primary" style="background: forestgreen;">Сохранить
        </button>
    </div>
</form>
{{template "footer.tpl"}}
</body>
</html>

