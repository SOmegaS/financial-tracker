import AuthForm from '../../common-components/auth-form';
import { useDispatch } from 'react-redux';
import { login } from '../../services/actions/login.ts';

export function SignInPage() {
    const dispatch = useDispatch();
    const handleSignIn = (email: string, password: string) => {
        // @ts-ignore
        dispatch(login(email, password));
    };

    return (
        <AuthForm
            title="Вход"
            buttonText="Войти"
            linkText="Еще не зарегистрированы?"
            linkTo="/sign-up"
            onSubmit={handleSignIn}
        />
    );
}
