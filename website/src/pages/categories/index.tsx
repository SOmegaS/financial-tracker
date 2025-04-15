import { Flex, Grid } from '@mantine/core';
import { initialState } from '../../services/initial-state.ts';
import { CategoryCard } from './components';
import { SectionHeader } from '../../common-components';

export const CategoriesPage = () => {
    return (
        <>
            <Flex direction={'column'}>
                <SectionHeader title={'Категории'} />
                <Grid gutter={{ base: 14 }}>
                    {initialState.categories.map((category) => (
                        <Grid.Col
                            key={category.id}
                            span={{ base: 12, sm: 6, md: 6, lg: 4 }}
                        >
                            <CategoryCard
                                name={category.name}
                                totalAmount={category.totalSum}
                            />
                        </Grid.Col>
                    ))}
                </Grid>
            </Flex>
        </>
    );
};
