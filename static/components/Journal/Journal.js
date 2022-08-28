class Journal {
    constructor() {
        this.url = '/api/v1/journal';
    }

    async fill_data() {
        // Описание таблицы
        let table_template = [
            { "head": "id", "row": "id" },
            { "head": "ID пользователя", "row": "user_id" },
            { "head": "Дата создания", "row": "created_at" },
            { "head": "Сумма", "row": "amount" }
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
            let table_users = document.querySelector("#table_journal");
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
        ROOT_JOURNAL.innerHTML = String.raw`    
            <hr>   
            <h3>Журнал транзакций</h3> 
            <table id="table_journal" class="table table-bordered"> 
                     <!-- Здесь будут сгенерированы элементы -->       
            </table>                      
        `;
        journal.fill_data();
    }
}

const journal = new Journal();