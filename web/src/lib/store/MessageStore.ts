export interface IMessageStore {
    addProgressMessage(message: string): void;
    removeProgressMessage(message: string): void;
}