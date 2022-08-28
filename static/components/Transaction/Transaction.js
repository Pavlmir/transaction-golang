class Transaction {
    constructor() {
        this.user = '';
        this.operation = 'Списание';
    }

    findUser(select) {
        this.user = select.querySelector(`option[value="${select.value}"]`).value;
    }

    findOption(select) {
        this.operation = select.querySelector(`option[value="${select.value}"]`).textContent;
    }

    async confirm_transaction(){
        let transaction_amount = document.querySelector("#transaction_amount");
         
        let amount = parseInt(transaction_amount.value)
        if (this.operation === "Списание") {
            amount = -amount;
        }
        
        let data = {
            'user': this.user,
            'amount': amount,
        }
    
        const response = await fetch(`/api/v1/transaction`, {
            method: 'POST',
            cache: 'no-cache',
            headers: {
                'Content-Type': 'application/json;charset=utf-8'
            },
            body: JSON.stringify(data),
        });

        const dataResponse = await response.text();
        alert(dataResponse);
        if (response.status !== 200) {
            errorPage.render(response);
        } else {
            userListPage.fill_data();
        }  
        journal.fill_data();
        
    }

    async fill_data() {
        // Структура шаблона
        // <option selected>Откройте это меню выбора</option>
        // <option value="1">Один</option>

        const response = await fetch(`/api/v1/users`, {
            method: 'GET'
        });

        if (response.status !== 200) {
            errorPage.render(response);
        } else {
            let user_list = await response.json();
            let select = document.querySelector("#transaction_users");

            for (let row of user_list) {
                let option = document.createElement("option");
                option.setAttribute("value", row.id);
                option.textContent = row.name;
                select.append(option);
            }
            this.user = select.querySelector("option").value;
        }
    }

    render() {
        let buttonСonfirm = String.raw`
        <button onclick="transaction.confirm_transaction();" class="btn btn-info">          
            <span> Подтвердить </span>
        </button>
        `;

        ROOT_TRANSACTION.innerHTML = String.raw` 
            <h3>Транзакция</h3>   
            <div class="transaction-block">
                <select id="transaction_users" class="form-select" aria-label="Выбор пользователей"
                        onchange='transaction.findUser(this)'>
                   <!-- Здесь будут сгенерированы элементы -->
                </select> 
                <select id="transaction_operation" class="form-select" aria-label="Выбор операции" 
                        onchange='transaction.findOption(this)'>
                    <option value="1">Списание</option>
                    <option value="2">Получение</option>
                </select>       
                <input id="transaction_amount" type="number" class="form-control" min="0">
                ${buttonСonfirm}
            </div>      
        `;

        transaction.fill_data();
    }
}

const transaction = new Transaction();