import { memo, useEffect, useState } from 'react';
import {
    Anchor,
    Divider,
    Group,
    Title,
    Popover,
    Button,
    ActionIcon,
} from '@mantine/core';
import styles from './header.module.css';
import { useLocation, useNavigate } from 'react-router-dom';
import { IconCoin, IconSettings, IconLogout } from '@tabler/icons-react';
import { useDispatch, useSelector } from 'react-redux';
import { IState } from '../../types.ts';
import { LOGOUT } from '../../services/actions';

function Header() {
    const navigate = useNavigate();
    const location = useLocation();
    const dispatch = useDispatch();
    const [active, setActive] = useState(0);
    const [popoverOpened, setPopoverOpened] = useState(false);

    const isAuth = useSelector((state: IState) => state.isAuth);

    useEffect(() => {
        if (location.pathname === '/sign-in') {
            setActive(1);
        } else if (location.pathname === '/sign-up') {
            setActive(2);
        } else {
            setActive(0);
        }
    }, [location.pathname]);

    const handleLogout = () => {
        dispatch({ type: LOGOUT });
        localStorage.removeItem('authToken');
        setPopoverOpened(false);
        navigate('/sign-in');
    };

    return (
        <header className={styles.header}>
            <Group align={'center'} justify="space-between">
                <Group
                    onClick={() => navigate('/categories')}
                    className={styles.title}
                    gap="xs"
                >
                    <Title order={1}>Монетка</Title>
                    <IconCoin
                        className={styles.icon}
                        size={40}
                        stroke={2}
                        color="var(--mantine-color-yellow-6)"
                    />
                </Group>

                <Group gap="md">
                    {isAuth ? (
                        <Popover
                            width={200}
                            position="bottom-end"
                            withArrow
                            shadow="md"
                            opened={popoverOpened}
                            onChange={setPopoverOpened}
                        >
                            <Popover.Target>
                                <ActionIcon
                                    variant="subtle"
                                    size="lg"
                                    aria-label="Settings"
                                    onClick={() => setPopoverOpened((o) => !o)}
                                >
                                    <IconSettings stroke={1.5} />
                                </ActionIcon>
                            </Popover.Target>
                            <Popover.Dropdown>
                                <Button
                                    fullWidth
                                    variant="light"
                                    color="red"
                                    leftSection={<IconLogout size={18} />}
                                    onClick={handleLogout}
                                >
                                    Выйти
                                </Button>
                            </Popover.Dropdown>
                        </Popover>
                    ) : null}

                    {!isAuth && (
                        <Group
                            gap={15}
                            visibleFrom="xm"
                            justify="flex-end"
                            className={styles.mainLinks}
                        >
                            <Anchor
                                key={'Вход'}
                                className={styles.mainLink}
                                data-active={1 === active || undefined}
                                onClick={() => {
                                    setActive(1);
                                    navigate('/sign-in');
                                }}
                            >
                                Вход
                            </Anchor>
                            <Anchor
                                key={'Регистрация'}
                                className={styles.mainLink}
                                data-active={2 === active || undefined}
                                onClick={() => {
                                    setActive(2);
                                    navigate('/sign-up');
                                }}
                            >
                                Регистрация
                            </Anchor>
                        </Group>
                    )}
                </Group>
            </Group>
            <Divider color="gray.4" size="xs" w="100%" />
        </header>
    );
}

export default memo(Header);
