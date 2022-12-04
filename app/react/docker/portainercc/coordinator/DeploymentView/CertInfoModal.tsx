import { useState } from 'react';
import { Database, Hash, Server, Tag, Tool } from 'react-feather';
import { DialogOverlay } from '@reach/dialog';

import { Button } from '@@/buttons';

import styles from 'react/sidebar/Footer/Footer.module.css'

export function CertInfoModalButton() {
  const [isCertInfoVisible, setIsCertInfoVisible] = useState(false);

  return (
    <>
      <button
        type="button"
        onClick={() => setIsCertInfoVisible(true)}
      >
        Inspect
      </button>
      {isCertInfoVisible && (
        <CertInfoModal closeModal={() => setIsCertInfoVisible(false)} />
      )}
    </>
  );
}

function CertInfoModal({ closeModal }: { closeModal: () => void }) {
  return (
    <DialogOverlay className={styles.dialog} isOpen>
      <div className="modal-dialog">
        <div className="modal-content">
          <div className="modal-header">
            <button type="button" className="close" onClick={closeModal}>
              ×
            </button>
            <h5 className="modal-title">Portainer</h5>
          </div>
          <div className="modal-body">
            <div className={styles.versionInfo}>
              <table>
                <tbody>
                  <tr>
                    <td>
                      <span className="inline-flex items-center">
                        <Server size="13" className="space-right" />
                        Server Version: 
                      </span>
                    </td>
                    <td>
                      <span className="inline-flex items-center">
                        <Database size="13" className="space-right" />
                        Database Version: 
                      </span>
                    </td>
                  </tr>
                  <tr>
                    <td>
                      <span className="inline-flex items-center">
                        <Hash size="13" className="space-right" />
                        CI Build Number:
                      </span>
                    </td>
                    <td>
                      <span>
                        <Tag size="13" className="space-right" />
                        Image Tag: 
                      </span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            <div className={styles.toolsList}>
              <span className="inline-flex items-center">
                <Tool size="13" className="space-right" />
                Compilation tools:
              </span>

              <div className={styles.tools}>
                <span className="text-muted small">
                  Nodejs v
                </span>
                <span className="text-muted small">
                  Yarn v
                </span>
                <span className="text-muted small">
                  Webpack
                </span>
                <span className="text-muted small">Go</span>
              </div>
            </div>
          </div>
          <div className="modal-footer">
            <Button className="bootbox-accept" onClick={closeModal}>
              Ok
            </Button>
          </div>
        </div>
      </div>
    </DialogOverlay>
  );
}
