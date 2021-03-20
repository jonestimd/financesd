import {RootStore} from './RootStore';
import * as entityUtils from '../model/entityUtils';
import * as agent from '../agent';
import {loadingPayees, query} from './PayeeStore';
import {newPayee, newPayeeModel} from 'src/test/payeeFactory';
import {PayeeModel} from '../model/PayeeModel';

describe('PayeeStore', () => {
    const {payeeStore, messageStore} = new RootStore();

    beforeEach(() => {
        payeeStore['payeesById'].clear();
    });
    describe('get payees', () => {
        it('sorts by name', () => {
            const payee = newPayeeModel();
            payeeStore['payeesById'].set(payee.id, payee);
            jest.spyOn(entityUtils, 'sortValuesByName');

            const payees = payeeStore.payees;

            expect(payees).toEqual([payee]);
            expect(entityUtils.sortValuesByName).toBeCalledWith(payeeStore['payeesById']);
        });
    });
    describe('getPayee', () => {
        it('returns payee for ID', () => {
            const payee = newPayeeModel();
            payeeStore['payeesById'].set(payee.id, payee);

            expect(payeeStore.getPayee(payee.id)).toBe(payee);
            expect(payeeStore.getPayee(parseInt(payee.id))).toBe(payee);
        });
        it('returns undefined for unknown ID', () => {
            expect(payeeStore.getPayee('-99')).toBeUndefined();
        });
    });
    describe('loadPayees', () => {
        beforeEach(() => {
            payeeStore['loading'] = false;
            jest.spyOn(messageStore, 'addProgressMessage');
            jest.spyOn(messageStore, 'removeProgressMessage');
        });
        it('loads payees if payeesById is empty', async () => {
            const payee = newPayee();
            jest.spyOn(agent, 'graphql').mockResolvedValue({data: {payees: [payee]}});

            await payeeStore.loadPayees();

            expect(payeeStore['loading']).toBe(false);
            expect(messageStore.addProgressMessage).toBeCalledWith(loadingPayees);
            expect(messageStore.removeProgressMessage).toBeCalledWith(loadingPayees);
            expect(agent.graphql).toBeCalledWith('/finances/api/v1/graphql', query);
            expect(payeeStore.payees).toStrictEqual([new PayeeModel(payee)]);
        });
        it('does nothing is already loading', async () => {
            payeeStore['loading'] = true;
            jest.spyOn(agent, 'graphql').mockRejectedValue(new Error());

            await payeeStore.loadPayees();

            expect(payeeStore['loading']).toBe(true);
            expect(messageStore.addProgressMessage).not.toBeCalled();
            expect(messageStore.removeProgressMessage).not.toBeCalled();
            expect(agent.graphql).not.toBeCalled();
        });
        it('does nothing is already loaded', async () => {
            const category = newPayeeModel();
            payeeStore['payeesById'].set(category.id, category);
            jest.spyOn(agent, 'graphql').mockRejectedValue(new Error());

            await payeeStore.loadPayees();

            expect(payeeStore['loading']).toBe(false);
            expect(messageStore.addProgressMessage).not.toBeCalled();
            expect(messageStore.removeProgressMessage).not.toBeCalled();
            expect(agent.graphql).not.toBeCalled();
        });
        it('logs error from graphql', async () => {
            const error = new Error('API error');
            jest.spyOn(agent, 'graphql').mockRejectedValue(error);
            jest.spyOn(console, 'error').mockImplementation(() => { });

            await payeeStore.loadPayees();

            expect(payeeStore['loading']).toBe(false);
            expect(messageStore.addProgressMessage).toBeCalledWith(loadingPayees);
            expect(messageStore.removeProgressMessage).toBeCalledWith(loadingPayees);
            expect(console.error).toBeCalledWith('error gettting payees', error);
        });
    });
});
