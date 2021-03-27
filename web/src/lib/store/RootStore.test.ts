import MessageStore from './MessageStore';
import {RootStore} from './RootStore';

describe('RootStore', () => {
    it('creates other stores', () => {
        const rootStore = new RootStore();

        expect(rootStore.messageStore).toBeInstanceOf(MessageStore);
        expect(rootStore.accountStore['loader']['messageStore']).toBe(rootStore.messageStore);
        expect(rootStore.categoryStore['loader']['messageStore']).toBe(rootStore.messageStore);
        expect(rootStore.groupStore['loader']['messageStore']).toBe(rootStore.messageStore);
        expect(rootStore.payeeStore['loader']['messageStore']).toBe(rootStore.messageStore);
        expect(rootStore.securityStore['loader']['messageStore']).toBe(rootStore.messageStore);
        expect(rootStore.transactionStore['rootStore']).toBe(rootStore);
    });
});
