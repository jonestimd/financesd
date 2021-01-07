import MessageStore from "./MessageStore";

describe('MessageStore', () => {
    describe('addProgressMessage', () => {
        it('appends message to list', () => {
            const messageStore = new MessageStore();

            messageStore.addProgressMessage('message 1');
            messageStore.addProgressMessage('message 2');

            expect(messageStore['progressMessages']).toEqual(['message 1', 'message 2']);
        });
    });
    describe('removeMessage', () => {
        it('removes first instance of message', () => {
            const messageStore = new MessageStore();
            messageStore.addProgressMessage('message 1');
            messageStore.addProgressMessage('message 2');
            messageStore.addProgressMessage('message 1');

            messageStore.removeProgressMessage('message 1');

            expect(messageStore['progressMessages']).toEqual(['message 2', 'message 1']);
        });
        it('does nothing if message not in list', () => {
            const messageStore = new MessageStore();
            messageStore.addProgressMessage('message 2');

            messageStore.removeProgressMessage('message 1');

            expect(messageStore['progressMessages']).toEqual(['message 2']);
        });
    });
    describe('get progressMessage', () => {
        it('returns undefined if no messages', () => {
            expect(new MessageStore().progressMessage).toBeUndefined();
        });
        it('returns first message', () => {
            const messageStore = new MessageStore();
            messageStore.addProgressMessage('message 1');
            messageStore.addProgressMessage('message 2');

            expect(messageStore.progressMessage).toEqual('message 1');

            messageStore.removeProgressMessage('message 1');
            expect(messageStore.progressMessage).toEqual('message 2');
        });
    });
});
