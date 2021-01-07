import MessageStore from "./MessageStore";
import {RootStore} from "./RootStore";

describe('RootStore', () => {
    it('creates other stores', () => {
        const rootStore = new RootStore();

        expect(rootStore.messageStore).toBeInstanceOf(MessageStore);
        expect(rootStore.accountStore['messageStore']).toBe(rootStore.messageStore);
        expect(rootStore.categoryStore['messageStore']).toBe(rootStore.messageStore);
        expect(rootStore.groupStore['messageStore']).toBe(rootStore.messageStore);
        expect(rootStore.payeeStore['messageStore']).toBe(rootStore.messageStore);
        expect(rootStore.securityStore['messageStore']).toBe(rootStore.messageStore);
        expect(rootStore.transactionStore['rootStore']).toBe(rootStore);
    });
});
