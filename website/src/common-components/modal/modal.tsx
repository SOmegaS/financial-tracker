import { Button, Modal, Text, Divider, Group } from '@mantine/core';
import React, { memo } from 'react';
import style from './modal.module.css';
import { IconCheck, IconX } from '@tabler/icons-react';

interface IModalProps {
    isOpen: boolean;
    onClose: () => void;
    title?: string;
    children?: React.ReactNode;
    onSave?: () => void;
    saveText?: string;
    cancelText?: string;
    size?: string | number;
    isDisabled?: boolean;
}

function CustomModal({
    title,
    isOpen,
    onClose,
    children,
    onSave = () => {},
    saveText = 'Сохранить',
    cancelText = 'Отмена',
    size = 'md',
    isDisabled = false,
}: IModalProps) {
    return (
        <Modal
            opened={isOpen}
            onClose={onClose}
            size={size}
            title={
                title ? (
                    <Text size="xl" fw={900}>
                        {title}
                    </Text>
                ) : null
            }
            overlayProps={{
                blur: 3,
                opacity: 0.55,
            }}
            transitionProps={{
                transition: 'fade',
                duration: 200,
            }}
            radius="md"
            padding="lg"
        >
            {children}

            <Divider my="lg" />

            <Group justify="flex-end" mt="md">
                <Button
                    variant="filled"
                    color="teal"
                    onClick={onSave}
                    disabled={isDisabled}
                    leftSection={<IconCheck className={style.icon} />}
                >
                    {saveText}
                </Button>
                <Button
                    variant="outline"
                    color="gray"
                    onClick={onClose}
                    leftSection={<IconX className={style.icon} />}
                >
                    {cancelText}
                </Button>
            </Group>
        </Modal>
    );
}

export default memo(CustomModal);
