import React from 'react';
import {shallow} from 'enzyme';
import TransactionList from './TransactionList';
import {RootStore} from 'src/lib/store/RootStore';
import TransactionTableModel from 'src/lib/model/TransactionTableModel';
import ListViewport from '../scroll/ListViewport';
import {newTxModel} from 'src/test/transactionFactory';
import Memo from './Memo';
import Checkbox from '@material-ui/core/Checkbox';
import Payee from './Payee';
import Security from './Security';
import TxDetail from './TxDetail';
import {newDetail} from 'src/test/detailFactory';

type RenderItem = Parameters<typeof ListViewport>[0]['renderItem'];

describe('TransactionList', () => {
    const {categoryStore, transactionStore} = new RootStore();
    const detail = newDetail({amount: 123.45});
    const txModel = newTxModel({categoryStore, memo: 'transaction memo', details: [detail], balance: 567.89});
    const transactionsModel = new TransactionTableModel([txModel], categoryStore);

    beforeEach(() => {
        jest.spyOn(React, 'useContext').mockReturnValue({transactionStore});
        jest.spyOn(transactionStore, 'getTransactionsModel').mockReturnValue(transactionsModel);
    });
    it('displays list in ListViewPort', () => {
        const component = shallow(<TransactionList />);

        expect(component.find(ListViewport).props()).toEqual(expect.objectContaining({
            items: transactionsModel.transactions,
            rowSelector: 'div.transaction',
            prototypeSelector: '.prototype',
        }));
    });
    it('displays prototype transaction', () => {
        const component = shallow(<TransactionList />);

        const prototype = component.find('TransactionPrototype').dive();
        expect(prototype.find('.leading .date')).toHaveText('0000-00-00');
        expect(prototype.find('.details').find(Memo)).toHaveProp('text', 'Prototype');
        expect(prototype.find('.trailing').find(Checkbox)).toExist();
        expect(prototype.find('.trailing .number').at(0)).toHaveText('0.00');
    });
    describe('renderItem', () => {
        let renderItem: RenderItem;

        beforeEach(() => {
            renderItem = shallow(<TransactionList />).find(ListViewport).prop('renderItem');
        });
        it('displays transaction date, payee, memo, security, details and cleared', () => {
            const tx = shallow(renderItem(txModel, undefined, false));

            expect(tx.find('.leading .date')).toHaveText(txModel.date);
            expect(tx.find(Payee)).toHaveProp('transaction', txModel);
            expect(tx.find(Security)).toHaveProp('transaction', txModel);
            expect(tx.find(Memo)).toHaveProp('text', txModel.memo);
            expect(tx.find('.details').find(TxDetail)).toHaveProp('detail', detail);
            expect(tx.find('.trailing').find(Checkbox)).toHaveProp('checked', false);
            expect(tx.find('.trailing .number').map((c) => c.text()))
                .toEqual([`$${txModel.subtotal}`, `$${txModel.balance}`]);
            expect(tx).not.toHaveClassName('selected');
            expect(tx.find('.ref-number')).not.toExist();
        });
        it('highlights selected transaction', () => {
            const tx = shallow(renderItem(newTxModel({categoryStore}), undefined, true));

            expect(tx).toHaveClassName('selected');
        });
        it('displays checked checkbox for cleared transaction', () => {
            const tx = shallow(renderItem(newTxModel({categoryStore, cleared: true}), undefined, true));

            expect(tx.find('.trailing').find(Checkbox)).toHaveProp('checked', true);
        });
        it('displays reference number', () => {
            const referenceNumber = '555';

            const tx = shallow(renderItem(newTxModel({referenceNumber, categoryStore}), undefined, true));

            expect(tx.find('.ref-number')).toHaveText(txModel.referenceNumber);
        });
    });
});
