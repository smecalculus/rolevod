package smecalculus.rolevod.core;

import static java.util.UUID.randomUUID;

import java.time.LocalDateTime;
import java.util.Collections;
import java.util.List;
import lombok.NonNull;
import lombok.RequiredArgsConstructor;
import smecalculus.rolevod.core.SepulkaMessageDm.PreviewRequest;
import smecalculus.rolevod.core.SepulkaMessageDm.PreviewResponse;
import smecalculus.rolevod.core.SepulkaMessageDm.RegistrationRequest;
import smecalculus.rolevod.core.SepulkaMessageDm.RegistrationResponse;
import smecalculus.rolevod.storage.SepulkaDao;

@RequiredArgsConstructor
public class SepulkaServiceImpl implements SepulkaService {

    @NonNull
    private SepulkaMapper mapper;

    @NonNull
    private SepulkaDao dao;

    @Override
    public RegistrationResponse register(RegistrationRequest request) {
        var now = LocalDateTime.now();
        var sepulkaCreated = mapper.toState(request)
                .internalId(randomUUID())
                .revision(0)
                .createdAt(now)
                .updatedAt(now)
                .build();
        var sepulkaSaved = dao.add(sepulkaCreated);
        return mapper.toMessage(sepulkaSaved).build();
    }

    @Override
    public List<PreviewResponse> view(PreviewRequest request) {
        return Collections.emptyList();
    }
}
