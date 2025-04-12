import {
    Anchor,
    Button,
    Container,
    Paper,
    PasswordInput,
    Text,
    TextInput,
    Title,
} from '@mantine/core';
import { useState } from 'react';
import { Link } from 'react-router-dom';

export function SignUpPage() {
    const [email, setEmail] = useState<string>('');
    const [password, setPassword] = useState<string>('');
    const [emailError, setEmailError] = useState<string | null>(null);
    const [passwordError, setPasswordError] = useState<string | null>(null);

    const validateEmail = (email: string) => {
        const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return re.test(email);
    };

    const handleEmailChange = (value: string) => {
        setEmail(value);
        if (value && !validateEmail(value)) {
            setEmailError('Некорректный формат почты');
        } else {
            setEmailError(null);
        }
    };

    const handlePasswordChange = (value: string) => {
        setPassword(value);
        if (value && value.length < 6) {
            setPasswordError('Пароль должен содержать минимум 6 символов');
        } else {
            setPasswordError(null);
        }
    };

    const handleSubmit = () => {
        if (!email) {
            setEmailError('Поле обязательно для заполнения');
            return;
        }
        if (!password) {
            setPasswordError('Поле обязательно для заполнения');
            return;
        }
        if (emailError || passwordError) {
            return;
        }

        console.log('Вход с:', { email, password });
    };

    return (
        <Container size="xs" my={100} w="100%" maw={400} px={20}>
            <Title ta="center" order={2}>
                Регистрация
            </Title>
            <Text c="dimmed" size="sm" ta="center" mt={5}>
                Уже зарегистрированы?
                <Anchor size="sm" component={Link} ml={5} to="/sign-in">
                    Войти
                </Anchor>
            </Text>

            <Paper withBorder shadow="md" p={30} mt={30} radius="md">
                <TextInput
                    onChange={(e) => handleEmailChange(e.target.value)}
                    value={email}
                    label="Почта"
                    placeholder="Ваша почта"
                    required
                    error={emailError}
                    withAsterisk
                />

                <PasswordInput
                    value={password}
                    onChange={(e) => handlePasswordChange(e.target.value)}
                    label="Пароль"
                    placeholder="Ваш пароль"
                    required
                    mt="md"
                    error={passwordError}
                    withAsterisk
                />

                <Button
                    fullWidth
                    mt="xl"
                    onClick={handleSubmit}
                    disabled={
                        !email || !password || !!emailError || !!passwordError
                    }
                >
                    Зарегистрироваться
                </Button>
            </Paper>
        </Container>
    );
}
