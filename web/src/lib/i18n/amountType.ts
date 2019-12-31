import {Bundle, defaultKey, translate} from './localize';

const bundle: Bundle = {
    [defaultKey]: {
        DEBIT_DEPOSIT: "Debit/Deposit",
        ASSET_VALUE: "Asset Value",
    }
};

export default function get(code: string) {
    return translate(code, bundle);
}