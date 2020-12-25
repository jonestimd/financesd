import {IBundle, defaultKey, translate} from './localize';

const bundle: IBundle = {
    [defaultKey]: {
        DEBIT_DEPOSIT: 'Debit/Deposit',
        ASSET_VALUE: 'Asset Value',
    },
};

export default function get(code: string) {
    return translate(code, bundle);
}
