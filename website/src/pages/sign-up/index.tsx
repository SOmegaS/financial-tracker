import AuthForm from '../../common-components/auth-form';
import { register } from '../../services/actions/register.ts';
import { useDispatch } from 'react-redux';

export function SignUpPage() {
    const dispatch = useDispatch();
    const handleSignUp = (email: string, password: string) => {
        dispatch(register(email, password));
    };

    return (
        <AuthForm
            title="Регистрация"
            buttonText="Зарегистрироваться"
            linkText="Уже зарегистрированы?"
            linkTo="/sign-in"
            onSubmit={handleSignUp}
        />
    );
}
