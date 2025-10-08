let token = localStorage.getItem('token');

// Проверка токена при загрузке страницы
if (token) {
    document.getElementById('auth-section').classList.add('d-none');
    document.getElementById('products-section').classList.remove('d-none');
    fetchProducts();
}

// Логин
document.getElementById('login-btn').addEventListener('click', async () => {
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    try {
        const response = await axios.post('http://localhost:8080/api/auth/login', {
            username,
            password
        });
        token = response.data.token;
        localStorage.setItem('token', token);

        // Скрываем форму логина, показываем продукты
        document.getElementById('auth-section').classList.add('d-none');
        document.getElementById('products-section').classList.remove('d-none');
        fetchProducts();
    } catch (error) {
        alert('Ошибка авторизации: ' + error.response?.data?.message || error.message);
    }
});

// Выход
document.getElementById('logout-btn').addEventListener('click', () => {
    localStorage.removeItem('token');
    token = null;
    document.getElementById('auth-section').classList.remove('d-none');
    document.getElementById('products-section').classList.add('d-none');
});

// Получение продуктов
async function fetchProducts() {
    try {
        const response = await axios.get('http://localhost:8080/api/products/all', {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        const products = response.data;
        const productsList = document.getElementById('products-list');
        productsList.innerHTML = products.map(product => `
            <tr>
                <td>${product.name}</td>
                <td>${product.quantity}</td>
                <td>${product.unit.abbreviation}</td>
            </tr>
        `).join('');
    } catch (error) {
        alert('Ошибка загрузки продуктов: ' + error.response?.data?.message || error.message);
    }
}
