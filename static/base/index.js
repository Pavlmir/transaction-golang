function render() {
    headerPage.render();  
    if (window.location.pathname === '/') {
        userListPage.render();
        transaction.render();
        journal.render();
    } 
}

render();
