import {newPayee} from 'test/payeeFactory';
import {PayeeModel} from './PayeeModel';

describe('PayeeModel', () => {
    const payee = newPayee();

    describe('constructor', () => {
        it('populates payee properties', () => {
            const model = new PayeeModel(payee);

            expect(model).toEqual(payee);
        });
    });
});
