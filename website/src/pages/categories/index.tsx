import { Flex, Title, Button, Group, Grid } from '@mantine/core';
import { initialState } from '../../services/initial-state.ts';
import { CategoryCard } from './components';
import { useState } from 'react';
import { AddCategoryModal } from './components';

export const CategoriesPage = () => {
    const [isOpenModal, setIsOpenModal] = useState(false);

    return (
        <>
            <Flex direction={'column'}>
                <Group align="center" justify="space-between">
                    <Title order={2}>Категории</Title>
                    <Button
                        variant="gradient"
                        gradient={{ from: 'teal', to: 'cyan', deg: 90 }}
                        size="compact-lg"
                        onClick={() => setIsOpenModal(true)}
                    >
                        Добавить категорию
                    </Button>
                </Group>
                <Grid gutter={{ base: 14 }}>
                    {initialState.categories.map((category) => (
                        <Grid.Col
                            key={category.id}
                            span={{ base: 12, sm: 6, md: 6, lg: 4 }}
                        >
                            <CategoryCard
                                id={category.id}
                                name={category.name}
                                totalAmount={category.totalSum}
                            />
                        </Grid.Col>
                    ))}
                </Grid>
            </Flex>
            <AddCategoryModal
                isOpen={isOpenModal}
                onClose={() => setIsOpenModal(false)}
            />
        </>
    );
};
