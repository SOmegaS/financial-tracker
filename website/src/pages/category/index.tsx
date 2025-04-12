import { Table, Text, Button, Paper, Group, Title, Badge } from '@mantine/core';
import { initialState } from '../../services/initial-state.ts';
import { IconTrash } from '@tabler/icons-react';
import { useState } from 'react';
import { AddReceiptModal, Modal } from '../../common-components';

function CategoryPage() {
    const [isAddReceiptOpenModal, setIsReceiptOpenModal] = useState(false);
    const [selectedReceipt, setSelectedReceipt] = useState<{
        name: string;
        id: string;
    } | null>(null);

    const handleDeleteClick = (receipt: { name: string; id: string }) => {
        setSelectedReceipt(receipt);
    };

    const handleConfirmDelete = () => {
        if (selectedReceipt) {
            console.log(`Удаляем чек: ${selectedReceipt.name}`);
        }
        setSelectedReceipt(null);
    };

    const rows = initialState.receipts.map((receipt) => (
        <Table.Tr key={receipt.id}>
            <Table.Td>
                <Text fw={500}>{receipt.name}</Text>
            </Table.Td>
            <Table.Td>
                {receipt.date.toLocaleDateString('ru-RU', {
                    day: '2-digit',
                    month: '2-digit',
                    year: 'numeric',
                })}
            </Table.Td>
            <Table.Td>
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
            <Table.Td style={{ display: 'flex', justifyContent: 'flex-end' }}>
                <Button
                    color="red.8"
                    leftSection={<IconTrash size={16} />}
                    onClick={() =>
                        handleDeleteClick({
                            name: receipt.name,
                            id: receipt.id,
                        })
                    }
                >
                    Удалить
                </Button>
            </Table.Td>
        </Table.Tr>
    ));

    return (
        <>
            <Paper p="md">
                <Group justify="space-between" mb="lg">
                    <Title order={2}>Чеки категории</Title>
                    <Button
                        variant="gradient"
                        gradient={{ from: 'teal', to: 'cyan', deg: 90 }}
                        size="compact-lg"
                        onClick={() => setIsReceiptOpenModal(true)}
                    >
                        Добавить чек
                    </Button>
                </Group>

                <Table striped highlightOnHover>
                    <Table.Thead>
                        <Table.Tr>
                            <Table.Th>Название</Table.Th>
                            <Table.Th>Дата</Table.Th>
                            <Table.Th>Сумма</Table.Th>
                            <Table.Th />
                        </Table.Tr>
                    </Table.Thead>
                    <Table.Tbody>{rows}</Table.Tbody>
                </Table>
            </Paper>

            <AddReceiptModal
                opened={isAddReceiptOpenModal}
                onClose={() => setIsReceiptOpenModal(false)}
                categoryName={'Здесь поменяй текст'}
            />

            <Modal
                title={`Вы уверены, что хотите удалить чек "${selectedReceipt?.name || ''}"?`}
                onClose={() => setSelectedReceipt(null)}
                isOpen={!!selectedReceipt}
                onSave={handleConfirmDelete}
                saveText="Удалить"
                cancelText="Отмена"
            />
        </>
    );
}

export default CategoryPage;
