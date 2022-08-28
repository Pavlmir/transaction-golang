class Header {
    render() {
        ROOT_HEADER.innerHTML = String.raw`    
            <a class="btn btn-secondary" id="docs" href="/docs" target="_blank">Документация</a> 
            <hr>                    
        `;
    }
}

const headerPage = new Header();