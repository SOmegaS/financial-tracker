import { Card, Text, Group, Button, Badge } from '@mantine/core';
import { IconReceipt, IconPlus } from '@tabler/icons-react';
import { memo, useState } from 'react';
import style from './category-card.module.css';
import { AddReceiptModal } from '../../../../common-components';
import { useNavigate } from 'react-router-dom';

interface CategoryCardProps {
    name: string;
    totalAmount: number;
    id: string;
}

function CategoryCard({ name, totalAmount, id }: CategoryCardProps) {
    const [isOpenModal, setIsOpenModal] = useState(false);
    const navigate = useNavigate();

    return (
        <>
            <Card
                withBorder
                shadow="sm"
                radius="md"
                p="lg"
                className={style.card}
            >
                <Group align="start" justify="space-between" mb="xs">
                    <Text fw={700} size="lg">
                        {name}
                    </Text>
                    <Badge size="lg" color="cyan.7">
                        <Text fw={700} size="lg">
                            {totalAmount.toLocaleString()} ₽
                        </Text>
                    </Badge>
                </Group>

                <Group mt="md" justify="space-between">
                    <Button
                        leftSection={<IconPlus className={style.icon} />}
                        variant="filled"
                        color="cyan.5"
                        onClick={() => setIsOpenModal(true)}
                    >
                        Добавить
                    </Button>
                    <Button
                        leftSection={<IconReceipt className={style.icon} />}
                        variant="light"
                        color="cyan.7"
                        onClick={() => navigate(`/category/${id}`)}
                    >
                        Показать все чеки
                    </Button>
                </Group>
            </Card>

            <AddReceiptModal
                categoryName={name}
                opened={isOpenModal}
                onClose={() => setIsOpenModal(false)}
            />
        </>
    );
}

export default memo(CategoryCard);
