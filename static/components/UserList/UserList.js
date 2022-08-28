class UserListPage {
    constructor() {
        this.url = '/api/v1/users';
    }

    async user_creation() {      
        let formElem = document.querySelector("#modal-form");
    
        const response = await fetch(`/api/v1/ctreate_user`, {
            method: 'POST',
            body: new FormData(formElem)
        });

        if (response.status !== 200) {
            errorPage.render(response);
        } 
        this.modal_windows();
    }

    modal_windows() {
        let modal_launcher = document.querySelector("#modal-launcher");
        let modal_background = document.querySelector("#modal-background");
        let modal_close = document.querySelector("#modal-close");
        let modal_content = document.querySelector("#modal-content");
        modal_launcher.classList.toggle("active");
        modal_background.classList.toggle("active");
        modal_close.classList.toggle("active");
        modal_content.classList.toggle("active");
    }

    user_update() {
        userListPage.fill_data();
    }

    async fill_data() {
        // Описание таблицы
        let table_template = [
            { "head": "id", "row": "id" },
            { "head": "Логин", "row": "name" },
            { "head": "Дата создания", "row": "created_at" },
            { "head": "Баланс", "row": "balance" }
        ]
        // Структура шаблона
        // <thead>
        //      <tr>
        //         <th></th>
        //      </th>
        // </thead>
        // <tbody>
        //      <tr>
        //         <td "data-value"=""></td>
        //      </tr>
        // </tbody>

        const response = await fetch(this.url, {
            method: 'GET'
        });

        if (response.status !== 200) {
            errorPage.render(response);
        } else {
            let user_list = await response.json();
            let table_users = document.querySelector("#table_users");
            table_users.innerHTML = '';

            let table_thead = document.createElement("thead");
            let tr_head = document.createElement("tr");
            table_thead.append(tr_head);
            table_users.append(table_thead);

            let table_tbody = document.createElement("tbody");
            let tr_body = document.createElement("tr");
            table_tbody.append(tr_body);
            table_users.append(table_tbody);

            for (let row_tmpl of table_template) {
                let th = document.createElement("th");
                th.innerHTML = row_tmpl.head;
                tr_head.append(th);

                let td = document.createElement("td");
                td.setAttribute("data-value", row_tmpl.row);
                tr_body.append(td);
            }

            let table_row = table_tbody.querySelector("tr");
            for (let row of user_list) {
                let tmpl_row = table_row.cloneNode(true);
                for (let td of tmpl_row.querySelectorAll("td")) {
                    td.innerHTML = row[td.dataset.value];
                }

                table_tbody.append(tmpl_row);
            }
            // Удаляем шаблон
            table_row.remove();
        }
    }

    render() {
        let buttonUserCreation = String.raw`
        <button onclick="userListPage.modal_windows();" class="btn btn-info"  id="modal-launcher">          
            <span> Создать нового пользователя </span>
       </button>
       `;

        let buttonUserUpdate = String.raw`
        <button onclick="userListPage.user_update();" class="btn btn-info">          
            <span> Обновить </span>
        </button>
        `;

        ROOT_USER_LIST.innerHTML = String.raw`   
            ${buttonUserCreation}
            ${buttonUserUpdate}
            <hr>
            <div>
                <h3>Список пользователей</h3>
                <table id="table_users" class="table table-bordered"> 
                     <!-- Здесь будут сгенерированы элементы -->       
                </table>       
            </div> 
            <div id="modal-background"></div>
            <div id="modal-content">
            <form id="modal-form" method="POST" action="">
                <div class="row mb-3">
                    <label for="inputName" class="col-sm-2 col-form-label">Имя</label>
                    <div class="col-sm-10">
                        <input type="text" class="form-control" name="name">
                    </div>
                </div>
                <br>
                <button type="submit" class="btn btn-primary" onclick="userListPage.user_creation();">Создать</button>
                <button id="modal-close" onclick="userListPage.modal_windows();">Отмена</button>
            </form>
            </div>
        `;

        userListPage.fill_data();
    }
}

const userListPage = new UserListPage();