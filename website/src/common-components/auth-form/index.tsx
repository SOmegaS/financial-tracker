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
import { memo, useState } from 'react';
import { Link } from 'react-router-dom';

interface AuthFormProps {
    title: string;
    buttonText: string;
    linkText: string;
    linkTo: string;
    onSubmit: (email: string, password: string) => void;
}

function AuthForm({
                      title,
                      buttonText,
                      linkText,
                      linkTo,
                      onSubmit,
                  }: AuthFormProps) {
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
        if (value && value.length <= 8) {
            setPasswordError('Пароль должен содержать минимум 9 символов');
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

        onSubmit(email, password);
    };

    return (
        <Container size="xs" my={100} w="100%" maw={400} px={20}>
            <Title ta="center" order={2}>
                {title}
            </Title>
            <Text c="dimmed" size="sm" ta="center" mt={5}>
                {linkText}
                <Anchor size="sm" component={Link} ml={5} to={linkTo}>
                    {title === 'Вход' ? 'Создать аккаунт' : 'Войти'}
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
                    {buttonText}
                </Button>
            </Paper>
        </Container>
    );
}

export default memo(AuthForm);
