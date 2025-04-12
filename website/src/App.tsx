import { Routes, Route, Navigate, BrowserRouter } from 'react-router-dom';
import style from './app.module.css';
import { CategoriesPage } from './pages/categories';
import { Header } from './../src/common-components';

function App() {
    return (
        <div className={style.main}>
            <BrowserRouter>
                <Header />
                <Routes>
                    <Route path={'/sign-in'} element={null} />
                    <Route path={'/sign-up'} element={null} />
                    <Route path={'/categories'} element={<CategoriesPage />} />
                    <Route
                        path="/"
                        element={<Navigate to="/categories" replace />}
                    />
                </Routes>
            </BrowserRouter>
        </div>
    );
}

export default App;
