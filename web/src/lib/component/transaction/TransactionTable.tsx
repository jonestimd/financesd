import React, {useCallback, useMemo} from 'react';
import {observer} from 'mobx-react-lite';
import {RootStoreContext} from '../../store/RootStore';
import {IColumn} from '../table/Row';
import HeaderDetailTable from '../table/HeaderDetailTable';
import TransactionModel, {ITransactionDetail} from '../../model/TransactionModel';
import * as formats from '../../formats';
import classNames from 'classnames';

interface IProps {
    accountId?: string;
}

const securityAccountTypes = ['BROKERAGE', '_401K'];

function numberClass(value: number, classes?: string) {
    return classNames(classes, 'number', {negative: value < 0});
}

const getDetails = (tx: TransactionModel) => tx.details;
const renderAmount = (detail: ITransactionDetail) => formats.currency.format(detail.amount);
const renderShares = (detail: ITransactionDetail) => {
    return detail.assetQuantity ? formats.shares.format(detail.assetQuantity) : '';
};
const dummyRender = () => '';

// TODO make header sticky
// TODO fix keyboard scrolling
// - slowness
// - always display selected row
const TransactionTable: React.FC<IProps> = observer(({accountId}) => {
    const {accountStore, categoryStore, groupStore, payeeStore, securityStore, transactionStore} = React.useContext(RootStoreContext);
    const account = accountStore.getAccount(accountId);
    const renderSecurity = useCallback((tx: TransactionModel) => securityStore.getSecurity(tx?.securityId)?.name ?? '', [securityStore]);
    // TODO filter security columns on non-security account
    const columns: IColumn<TransactionModel>[] = useMemo(() => [
        {key: 'transaction.date', render: (tx) => tx.date, className: 'date'},
        {key: 'transaction.referenceNumber', render: (tx) => tx.referenceNumber},
        {key: 'transaction.payee', render: (tx) => payeeStore.getPayee(tx?.payeeId)?.name ?? ''},
        {key: 'transaction.memo', render: (tx) => tx.memo},
        {key: 'transaction.security', render: renderSecurity, className: 'security'},
        {key: 'transaction.subtotal', render: (tx) => formats.currency.format(tx.subtotal), className: (tx) => numberClass(tx?.subtotal ?? 0)},
        {key: 'transaction.cleared', render: (tx) => tx.cleared ? <span>&#x2713;</span> : null, className: 'boolean'},
        {key: 'transaction.balance', render: (tx) => formats.currency.format(tx.balance), className: (tx) => numberClass(tx?.balance ?? 0)},
    ], [payeeStore, renderSecurity]);
    const renderCategory = useCallback((detail: ITransactionDetail) => {
        if (detail.relatedDetail) {
            return <span className='transfer'>{accountStore.getAccount(detail.relatedDetail.transaction.accountId)?.name}</span>;
        }
        return <span>{categoryStore.getCategory(detail.transactionCategoryId)?.displayName ?? ''}</span>;
    }, [accountStore, categoryStore]);
    const renderGroup = useCallback(({transactionGroupId}: ITransactionDetail) => {
        return <span className='group'>{groupStore.getGroup(transactionGroupId)?.name ?? ''}</span>;
    }, [groupStore]);
    const subcolumns: IColumn<ITransactionDetail>[] = useMemo(() => [
        {key: 'detail.group', colspan: 2, render: renderGroup, className: 'group'},
        {key: 'detail.category', render: renderCategory, className: 'category'},
        {key: 'detail.memo', render: (detail) => detail.memo},
        {key: 'detail.shares', render: renderShares, className: (detail) => numberClass(detail?.assetQuantity ?? 0, 'security')},
        {key: 'detail.amount', render: renderAmount, className: (detail) => numberClass(detail?.amount ?? 0)},
        {key: 'dummy1', header: dummyRender, render: dummyRender},
        {key: 'dummy2', header: dummyRender, render: dummyRender},
    ], [renderCategory, renderGroup]);
    return (
        <HeaderDetailTable
            className={securityAccountTypes.includes(account?.type ?? '') ? 'security-transactions' : 'transactions'}
            columns={columns}
            subColumns={subcolumns}
            model={transactionStore.getTransactionsModel(accountId)}
            subrows={getDetails} />
    );
});

export default TransactionTable;
