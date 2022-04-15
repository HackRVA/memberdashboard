interface HTMLElementTagNameMap {
    'edit-member-form': Element & WithProperties<{
        closeHandler: Function;
    }>;
    'abstract-dialog': Dialog & WithProperties<{
        id: string;
        heading: string;
        dialogLayout: TemplateResult;
        handleClose: Function;
        handleSubmit: Function;
        toastMessage: string;
        dialogOpened: boolean;
    }>;
    'member-grid': Element & WithProperties<{
        members: Member[];
    }>
    
    'abstract-toast': Element & WithProperties<{
        message: string;
        timeoutMs: number;
    }>
    'toast-message': Element & WithProperties<{
        message: string;
        timeoutMs: number;
    }>
  }