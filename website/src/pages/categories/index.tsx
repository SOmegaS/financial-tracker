import { Flex, Grid } from '@mantine/core';
import { CategoryCard } from './components';
import { SectionHeader } from '../../common-components';
import { useDispatch, useSelector } from 'react-redux';
import { memo, useEffect } from 'react';
import { IState } from '../../types.ts';
import { getCategories } from '../../services/actions/get-categories.ts';

const CategoriesPage = () => {
    const dispatch = useDispatch();
    const categories = useSelector((state: IState) => state.categories);
    const isLoading = useSelector((state: IState) => state.isLoading);
    useEffect(() => {
        if (!isLoading) {
            // @ts-ignore
            dispatch(getCategories());
        }
    }, [dispatch]);

    return (
        <>
            <Flex direction={'column'}>
                <SectionHeader title={'Категории'} />
                <Grid gutter={{ base: 14 }}>
                    {categories &&
                        categories.map((category) => (
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

export default memo(CategoriesPage);
