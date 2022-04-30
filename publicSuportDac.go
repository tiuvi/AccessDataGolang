package bd

//Verifica si el string, es una columna en ese espacio.
func (SF *spaceFile) IsColumn(column string) *bool {

	if SF.indexSizeColumns != nil {

		_, found := SF.indexSizeColumns[column]
		if found {

			exist := true
			return &exist
		}

		exist := false
		return &exist
	}

	return nil
}

//Verifica si el string, es un campo en ese espacio.
func (SF *spaceFile) IsField(field string) *bool {

	if SF.indexSizeFields != nil {

		_, found := SF.indexSizeFields[field]
		if found {

			exist := true
			return &exist
		}

		exist := false
		return &exist
	}

	return nil
}