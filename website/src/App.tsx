import { Routes, Route, Navigate, BrowserRouter } from 'react-router-dom';
import style from './app.module.css';
import CategoriesPage from './pages/categories';
import { Header } from './../src/common-components';
import CategoryPage from './pages/category';
import { SignInPage } from './pages/sign-in';
import { SignUpPage } from './pages/sign-up';
import { useDispatch, useSelector } from 'react-redux';
import { LOGIN, SET_ERROR } from './services/actions';
import { notifications } from '@mantine/notifications';
import { useEffect } from 'react';
import { IState } from './types.ts';
import { LoadingOverlay } from '@mantine/core';

function App() {
    const dispatch = useDispatch();
    const error = useSelector((state: IState) => state.error);
    const isAuth = useSelector((state: IState) => state.isAuth);
    const isLoading = useSelector((state: IState) => state.isLoading);
    useEffect(() => {
        if (error) {
            notifications.show({
                title: 'Ошибка',
                message: error,
                position: 'top-right',
                color: 'red',
                autoClose: 5000,
                onClose: () => dispatch({ type: SET_ERROR, payload: null }),
            });

            dispatch({ type: SET_ERROR, payload: null });
        }
    }, [error, dispatch]);

    useEffect(() => {
        const token = localStorage.getItem('authToken');
        if (token) {
            dispatch({
                type: LOGIN,
                payload: token,
            });
        }
    }, []);

    return (
        <div className={style.main}>
            <BrowserRouter>
                <Header />
                <LoadingOverlay
                    visible={isLoading}
                    overlayProps={{ blur: 2 }}
                    loaderProps={{
                        type: 'bars',
                        color: 'blue',
                        size: 'xl',
                    }}
                    zIndex={1000}
                />
                <Routes>
                    {isAuth ? (
                        <>
                            <Route
                                path={'/categories'}
                                element={<CategoriesPage />}
                            />
                            <Route
                                path={'/category/:name'}
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
