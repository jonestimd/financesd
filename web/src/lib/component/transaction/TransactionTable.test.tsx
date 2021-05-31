import React from 'react';
import {shallow, ShallowWrapper} from 'enzyme';
import TransactionTable from './TransactionTable';
import {RootStore} from 'lib/store/RootStore';
import {newAccountModel} from 'test/accountFactory';
import HeaderDetailTable from '../table/HeaderDetailTable';
import TransactionTableModel from 'lib/model/TransactionTableModel';
import {newTxModel} from 'test/transactionFactory';
import {newPayeeModel} from 'test/payeeFactory';
import {newSecurityModel} from 'test/securityFactory';
import {newDetail} from 'test/detailFactory';
import {newGroupModel} from 'test/groupFactory';
import {newCategoryModel} from 'test/categoryFactory';
import TransactionModel, {ITransactionDetail} from 'lib/model/TransactionModel';

const account = newAccountModel();

const getTxColumns = (component: ShallowWrapper) => component.find(HeaderDetailTable).prop('columns');
const getTxColumn = (component: ShallowWrapper, key: string) => getTxColumns(component).find((col) => col.key === key);
const getDetailColumns = (component: ShallowWrapper) => component.find(HeaderDetailTable).prop('subColumns');
const getDetailColumn = (component: ShallowWrapper, key: string) => getDetailColumns(component).find((col) => col.key === key);

describe('TransactionTable', () => {
    const rootStore = new RootStore();
    const transactionsModel = new TransactionTableModel([], rootStore.categoryStore);
    const payee = newPayeeModel();
    const security = newSecurityModel();
    const txModel = newTxModel({
        payeeId: payee.id,
        securityId: security.id,
        referenceNumber: '123',
        memo: 'tx memo',
        details: [newDetail()],
        balance: 987.23,
    });

    beforeEach(() => {
        jest.spyOn(React, 'useContext').mockReturnValue(rootStore);
        jest.spyOn(rootStore.accountStore, 'getAccount').mockReturnValueOnce(account);
        jest.spyOn(rootStore.transactionStore, 'getTransactionsModel').mockReturnValue(transactionsModel);
        jest.spyOn(rootStore.payeeStore, 'getPayee').mockReturnValue(payee);
        rootStore.securityStore['securitiesById'].set(security.id, security);
    });
    it('displays transactions in a header/detail table', () => {
        const component = shallow(<TransactionTable accountId={account.id} />);

        expect(component.find(HeaderDetailTable)).toEqual(component);
        expect(component).toHaveClassName('transactions');
        expect(component).toHaveProp('model', transactionsModel);
        expect(component.find(HeaderDetailTable).prop('subrows')(txModel)).toEqual(txModel.details);
    });
    describe('transaction row', () => {
        it('displays date', () => {
            const component = shallow(<TransactionTable accountId={account.id} />);
            const renderer = getTxColumn(component, 'transaction.date')!.render;

            expect(renderer(txModel)).toEqual(txModel.date);
        });
        it('displays reference number', () => {
            const component = shallow(<TransactionTable accountId={account.id} />);
            const renderer = getTxColumn(component, 'transaction.referenceNumber')!.render;

            expect(renderer(txModel)).toEqual(txModel.referenceNumber);
        });
        it('displays payee', () => {
            const component = shallow(<TransactionTable accountId={account.id} />);
            const renderer = getTxColumn(component, 'transaction.payee')!.render;

            expect(renderer(txModel)).toEqual(payee.name);

            expect(rootStore.payeeStore.getPayee).toBeCalledWith(txModel.payeeId);
        });
        it('displays blank for no security', () => {
            const component = shallow(<TransactionTable accountId={account.id} />);
            const renderer = getTxColumn(component, 'transaction.security')!.render;

            expect(renderer(newTxModel())).toEqual('');
        });
        it('displays security', () => {
            const component = shallow(<TransactionTable accountId={account.id} />);
            const renderer = getTxColumn(component, 'transaction.security')!.render;

            expect(renderer(txModel)).toEqual(security.name);
        });
        it('displays memo', () => {
            const component = shallow(<TransactionTable accountId={account.id} />);
            const renderer = getTxColumn(component, 'transaction.memo')!.render;

            expect(renderer(txModel)).toEqual(txModel.memo);
        });
        it('displays subtotal', () => {
            const component = shallow(<TransactionTable accountId={account.id} />);
            const {render: renderer, ...rest} = getTxColumn(component, 'transaction.subtotal')!;

            expect(renderer(txModel)).toEqual(`$${txModel.subtotal}`);
            const className = rest.className as (tx: Partial<TransactionModel>) => string;
            expect(className({subtotal: 0})).toEqual('number');
            expect(className({subtotal: -1})).toEqual('number negative');
        });
        it('displays balance', () => {
            const component = shallow(<TransactionTable accountId={account.id} />);
            const {render: renderer, ...rest} = getTxColumn(component, 'transaction.balance')!;

            expect(renderer(txModel)).toEqual(`$${txModel.balance}`);
            const className = rest.className as (tx: Partial<TransactionModel>) => string;
            expect(className({balance: 0})).toEqual('number');
            expect(className({balance: -1})).toEqual('number negative');
        });
        it('displays cleared status', () => {
            const component = shallow(<TransactionTable accountId={account.id} />);
            const renderer = getTxColumn(component, 'transaction.cleared')!.render;

            expect(renderer(txModel)).toBeNull();
            expect(renderer({cleared: true})).toEqual(<span>&#x2713;</span>);
        });
    });
    describe('detail row', () => {
        const group = newGroupModel();
        const category = newCategoryModel();
        const detail = newDetail({
            transactionGroupId: group.id,
            transactionCategoryId: category.id,
            memo: 'detail memo',
            assetQuantity: 9.234567,
            amount: 456.23,
        });

        it('displays blank for no group', () => {
            const component = shallow(<TransactionTable accountId={account.id} />);
            const renderer = getDetailColumn(component, 'detail.group')!.render;

            expect(renderer(newDetail())).toEqual(<span className='group'>{''}</span>);
        });
        it('displays group', () => {
            rootStore.groupStore['groupsById'].set(group.id, group);
            const component = shallow(<TransactionTable accountId={account.id} />);
            const renderer = getDetailColumn(component, 'detail.group')!.render;

            expect(renderer(detail)).toEqual(<span className='group'>{group.name}</span>);
        });
        it('displays blank for no category', () => {
            const component = shallow(<TransactionTable accountId={account.id} />);
            const renderer = getDetailColumn(component, 'detail.category')!.render;

            expect(renderer(newDetail())).toEqual(<span>{''}</span>);
        });
        it('displays category', () => {
            rootStore.categoryStore['categoriesById'].set(category.id, category);
            const component = shallow(<TransactionTable accountId={account.id} />);
            const renderer = getDetailColumn(component, 'detail.category')!.render;

            expect(renderer(detail)).toEqual(<span>{category.displayName}</span>);
        });
        it('displays transfer account', () => {
            const relatedAccount = newAccountModel();
            const relatedTx = {accountId: relatedAccount.id};
            const relatedDetail = {transaction: relatedTx};
            jest.spyOn(rootStore.accountStore, 'getAccount').mockReturnValueOnce(relatedAccount);
            const component = shallow(<TransactionTable accountId={account.id} />);
            const renderer = getDetailColumn(component, 'detail.category')!.render;

            expect(renderer({relatedDetail})).toEqual(<span className='transfer'>{relatedAccount.name}</span>);

            expect(rootStore.accountStore.getAccount).toBeCalledWith(relatedTx.accountId);
        });
        it('displays memo', () => {
            const component = shallow(<TransactionTable accountId={account.id} />);
            const renderer = getDetailColumn(component, 'detail.memo')!.render;

            expect(renderer(detail)).toEqual(detail.memo);
        });
        it('displays shares', () => {
            const component = shallow(<TransactionTable accountId={account.id} />);
            const {render: renderer, ...rest} = getDetailColumn(component, 'detail.shares')!;

            // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
            expect(renderer(detail)).toEqual(`${detail.assetQuantity}`);
            expect(renderer({})).toEqual('');
            const className = rest.className as (detail: Partial<ITransactionDetail>) => string;
            expect(className({assetQuantity: 0})).toEqual('security number');
            expect(className({assetQuantity: -1})).toEqual('security number negative');
        });
        it('displays amount', () => {
            const component = shallow(<TransactionTable accountId={account.id} />);
            const {render: renderer, ...rest} = getDetailColumn(component, 'detail.amount')!;

            expect(renderer(detail)).toEqual(`$${detail.amount}`);
            const className = rest.className as (detail: Partial<ITransactionDetail>) => string;
            expect(className({amount: 0})).toEqual('number');
            expect(className({amount: -1})).toEqual('number negative');
        });
        it('displays filler columns', () => {
            const component = shallow(<TransactionTable accountId={account.id} />);

            expect(getDetailColumn(component, 'dummy1')?.render({})).toEqual('');
            expect(getDetailColumn(component, 'dummy2')?.render({})).toEqual('');
        });
    });
});
