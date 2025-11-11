# vendor.proto.mk для Windows

VENDOR_PROTO_PATH := vendor.protobuf

vendor: .vendor-reset .vendor-googleapis .vendor-google-protobuf .vendor-protovalidate .vendor-protoc-gen-openapiv2 .vendor-tidy

.vendor-reset:
	if exist "$(VENDOR_PROTO_PATH)" rmdir /s /q "$(VENDOR_PROTO_PATH)"
	mkdir "$(VENDOR_PROTO_PATH)"

.vendor-protovalidate:
	git clone -b main --single-branch --depth=1 --filter=tree:0 https://github.com/bufbuild/protovalidate "$(VENDOR_PROTO_PATH)/protovalidate"
	cd "$(VENDOR_PROTO_PATH)/protovalidate" && git checkout
	if exist "$(VENDOR_PROTO_PATH)\protovalidate\proto\protovalidate\buf" move "$(VENDOR_PROTO_PATH)\protovalidate\proto\protovalidate\buf" "$(VENDOR_PROTO_PATH)\"
	if exist "$(VENDOR_PROTO_PATH)/protovalidate" rmdir /s /q "$(VENDOR_PROTO_PATH)/protovalidate"

.vendor-google-protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 https://github.com/protocolbuffers/protobuf "$(VENDOR_PROTO_PATH)/protobuf"
	cd "$(VENDOR_PROTO_PATH)/protobuf" && git sparse-checkout set --no-cone src/google/protobuf && git checkout
	if not exist "$(VENDOR_PROTO_PATH)/google/protobuf" mkdir "$(VENDOR_PROTO_PATH)/google/protobuf"
	if exist "$(VENDOR_PROTO_PATH)\protobuf\src\google\protobuf\*.proto" move "$(VENDOR_PROTO_PATH)\protobuf\src\google\protobuf\*.proto" "$(VENDOR_PROTO_PATH)\google\protobuf\"
	if exist "$(VENDOR_PROTO_PATH)/protobuf" rmdir /s /q "$(VENDOR_PROTO_PATH)/protobuf"

.vendor-googleapis:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 https://github.com/googleapis/googleapis "$(VENDOR_PROTO_PATH)/googleapis"
	cd "$(VENDOR_PROTO_PATH)/googleapis" && git sparse-checkout set --no-cone google/api google/type && git checkout
	if not exist "$(VENDOR_PROTO_PATH)/google/api" mkdir "$(VENDOR_PROTO_PATH)/google/api"
	if not exist "$(VENDOR_PROTO_PATH)/google/type" mkdir "$(VENDOR_PROTO_PATH)/google/type"
	if exist "$(VENDOR_PROTO_PATH)\googleapis\google\api\*.proto" move "$(VENDOR_PROTO_PATH)\googleapis\google\api\*.proto" "$(VENDOR_PROTO_PATH)\google\api\"
	if exist "$(VENDOR_PROTO_PATH)\googleapis\google\type\*.proto" move "$(VENDOR_PROTO_PATH)\googleapis\google\type\*.proto" "$(VENDOR_PROTO_PATH)\google\type\"
	if exist "$(VENDOR_PROTO_PATH)/googleapis" rmdir /s /q "$(VENDOR_PROTO_PATH)/googleapis"

.vendor-protoc-gen-openapiv2:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 https://github.com/grpc-ecosystem/grpc-gateway "$(VENDOR_PROTO_PATH)/grpc-gateway"
	cd "$(VENDOR_PROTO_PATH)/grpc-gateway" && git sparse-checkout set --no-cone protoc-gen-openapiv2/options && git checkout
	if not exist "$(VENDOR_PROTO_PATH)/protoc-gen-openapiv2" mkdir "$(VENDOR_PROTO_PATH)/protoc-gen-openapiv2"
	if exist "$(VENDOR_PROTO_PATH)\grpc-gateway\protoc-gen-openapiv2\options" move "$(VENDOR_PROTO_PATH)\grpc-gateway\protoc-gen-openapiv2\options" "$(VENDOR_PROTO_PATH)\protoc-gen-openapiv2\"
	if exist "$(VENDOR_PROTO_PATH)/grpc-gateway" rmdir /s /q "$(VENDOR_PROTO_PATH)/grpc-gateway"

.vendor-tidy:
	powershell -Command "Get-ChildItem '$(VENDOR_PROTO_PATH)' -Recurse -File | Where-Object { $$_.Extension -ne '.proto' } | Remove-Item -Force"
	powershell -Command "Get-ChildItem '$(VENDOR_PROTO_PATH)' -Recurse -File | Where-Object { $$_.Name -match 'unittest|test|sample' } | Remove-Item -Force"
	powershell -Command "Get-ChildItem '$(VENDOR_PROTO_PATH)' -Recurse -Directory | Where-Object { (Get-ChildItem $$_.FullName).Count -eq 0 } | Remove-Item -Force"

.PHONY: .vendor-reset .vendor-google-protobuf .vendor-googleapis .vendor-protoc-gen-openapiv2 .vendor-protovalidate .vendor-tidy vendor