import { Button, Group, Title } from '@mantine/core';
import { memo, useState } from 'react';
import { AddReceiptModal } from '../add-receipt-modal';
import { useParams } from 'react-router-dom';

function SectionHeader({ title }: { title: string }) {
    const { name } = useParams();
    const decodedName = name ? decodeURIComponent(atob(name)) : null;

    const [isAddReceiptOpenModal, setIsReceiptOpenModal] = useState(false);
    return (
        <>
            <Group justify="space-between" mb="lg">
                <Title order={2}>{title}</Title>
                <Button
                    variant="gradient"
                    gradient={{ from: 'teal', to: 'cyan', deg: 90 }}
                    size="compact-lg"
                    onClick={() => setIsReceiptOpenModal(true)}
                >
                    Добавить чек
                </Button>
            </Group>
            <AddReceiptModal
                categoryName={decodedName || ''}
                opened={isAddReceiptOpenModal}
                onClose={() => setIsReceiptOpenModal(false)}
            />
        </>
    );
}

export default memo(SectionHeader);
