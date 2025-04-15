import AuthForm from '../../common-components/auth-form';

export function SignInPage() {
    const handleSignIn = (email: string, password: string) => {
        console.log('Вход с:', { email, password });
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
