import React from 'react';
import {shallow} from 'enzyme';
import TransactionModel from 'lib/model/TransactionModel';
import {PayeeModel} from 'lib/model/PayeeModel';
import Payee from './Payee';
import {newTx} from 'test/transactionFactory';
import {RootStore} from 'lib/store/RootStore';
import {Icon} from '@material-ui/core';

const txData = newTx();

describe('Payee', () => {
    const rootStore = new RootStore();
    const {payeeStore, accountStore, categoryStore} = rootStore;

    beforeEach(() => {
        jest.spyOn(React, 'useContext').mockReturnValue({payeeStore});
    });
    it('returns null if no payee', () => {
        const tx = new TransactionModel(txData, accountStore, categoryStore);

        const component = shallow(<Payee transaction={tx} />);

        expect(component).toBeEmptyRender();
    });
    it('shows payee', () => {
        const name = 'the payee';
        jest.spyOn(payeeStore, 'getPayee').mockReturnValue({name} as PayeeModel);
        const tx = new TransactionModel({...txData, payeeId: -1}, accountStore, categoryStore);

        const component = shallow(<Payee transaction={tx} />);

        expect(component).toHaveProp('data-type', 'payee');
        expect(component).toHaveClassName('chip');
        expect(component.find(Icon)).toHaveText('person');
        expect(component.childAt(1)).toHaveText(name);
    });
});
