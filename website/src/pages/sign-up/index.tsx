import AuthForm from '../../common-components/auth-form';

export function SignUpPage() {
    const handleSignUp = (email: string, password: string) => {
        console.log('Регистрация с:', { email, password });
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
