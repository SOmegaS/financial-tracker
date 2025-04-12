import { Routes, Route, Navigate, BrowserRouter } from 'react-router-dom';
import style from './app.module.css';
import { CategoriesPage } from './pages/categories';
import { Header } from './../src/common-components';
import CategoryPage from './pages/category';
import { SignInPage } from './pages/sign-in';
import { SignUpPage } from './pages/sign-up';
import { initialState } from './services/initial-state.ts';

function App() {
    return (
        <div className={style.main}>
            <BrowserRouter>
                <Header />
                <Routes>
                    {initialState.isLogin ? (
                        <>
                            <Route
                                path={'/categories'}
                                element={<CategoriesPage />}
                            />
                            <Route
                                path={'/category/:id'}
                                element={<CategoryPage />}
                            />
                            <Route
                                path="*"
                                element={<Navigate to="/categories" replace />}
                            />
                        </>
                    ) : (
                        <>
                            <Route path={'/sign-in'} element={<SignInPage />} />
                            <Route path={'/sign-up'} element={<SignUpPage />} />
                            <Route
                                path="*"
                                element={<Navigate to="/sign-in" replace />}
                            />
                        </>
                    )}
                </Routes>
            </BrowserRouter>
        </div>
    );
}

export default App;
