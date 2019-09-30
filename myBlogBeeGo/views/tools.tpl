<div>
    <a uk-tooltip="title: Редактировать; pos: top" class="uk-icon-button "
       href="/posts/edit/?id={{.ID}}" uk-icon="icon: file-edit" style="color: deepskyblue;"></a>
    <a onclick="
            if (confirm('Вы уверены, что хотите УДАЛИТЬ этот пост?')) {
            deletePost('{{.ID}}')}
            else {
            return false
            }" uk-tooltip="title: Удалить; pos: top" class="uk-icon-button" href="#" uk-icon="icon: trash"
       style="color: red;"></a>
</div>
