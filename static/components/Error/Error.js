class ErrorPage {
    render(response) {
        let statusText = response.statusText
        let status = response.status
        ROOT_ERROR.innerHTML = String.raw`               
             <div class="cover-container d-flex w-90 h-100 p-3 mx-auto flex-column bg-light p-5 rounded mt-4">                
                <main class="px-6" style="margin-left: 15%">
                    <h1>Page not found.</h1>
                    <p class="lead">The request URL was not found on the server. If you entered the URL manually please
                        check
                        your spelling and try again.</p>
                    <p class="lead"> ${statusText} - ${status} </p>
                    <p class="lead">
                        <a href="{{ url_for('login_page') }}"
                        class="btn btn-lg btn-secondary fw-bold border-white bg-white text-dark">Back to main page</a>
                    </p>
                </main>
             </div>         
        `;
    }
}

const errorPage = new ErrorPage();