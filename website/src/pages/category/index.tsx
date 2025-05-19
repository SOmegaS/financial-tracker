import { Table, Text, Paper, Badge } from '@mantine/core';
import { SectionHeader } from '../../common-components';
import style from './category.module.css';
import { useParams } from 'react-router-dom';
import { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { IState } from '../../types.ts';
import { getBills } from '../../services/actions/get-biils.ts';

function CategoryPage() {
    const { name } = useParams();
    const decodedName = name ? decodeURIComponent(atob(name)) : '';
    const receipts = useSelector((state: IState) => state.receipts);
    const dispatch = useDispatch();
    useEffect(() => {
        // @ts-ignore
        dispatch(getBills(decodedName));
    }, [dispatch]);

    if (!receipts || !receipts.length) {
        return null;
    }

    const rows = receipts.map((receipt) => (
        <Table.Tr key={receipt.id}>
            <Table.Td>
                <Text fw={500}>{receipt.name}</Text>
            </Table.Td>
            <Table.Td>{receipt.date}</Table.Td>
            <Table.Td className={style.sum}>
                <Badge
                    size="lg"
                    variant="gradient"
                    gradient={{ from: 'cyan.8', to: 'cyan.6', deg: 25 }}
                >
                    <Text fw={700} size="md">
                        {receipt.sum} ₽
                    </Text>
                </Badge>
            </Table.Td>
        </Table.Tr>
    ));

    return (
        <Paper p="md">
            <SectionHeader title={`Чеки категории ${decodedName}`} />

            <Table striped highlightOnHover>
                <Table.Thead>
                    <Table.Tr>
                        <Table.Th>Название</Table.Th>
                        <Table.Th>Дата</Table.Th>
                        <Table.Th className={style.sum}>Сумма</Table.Th>
                    </Table.Tr>
                </Table.Thead>
                <Table.Tbody>{rows}</Table.Tbody>
            </Table>
        </Paper>
    );
}

export default CategoryPage;
